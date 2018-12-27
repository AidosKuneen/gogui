// Copyright (c) 2019 Aidos Developer

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package gogui

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

//Message Headers
const (
	headerConnect = iota
	headerEvent
	headerACK
	headerError
)

type packetRead struct {
	Type  byte            `json:"type"`
	ID    uint            `json:"id"`
	Param json.RawMessage `json:"param"`
}

type packetWrite struct {
	Type  byte        `json:"type"`
	ID    uint        `json:"id"`
	Param interface{} `json:"param"`
}

type eventParamRead struct {
	Name  string          `json:"name"`
	Param json.RawMessage `json:"param"`
}

type eventParamWrite struct {
	Name  string      `json:"name"`
	Param interface{} `json:"param"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 65535
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	sync.RWMutex
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *packetWrite

	fconnect    func()
	fdisconnect func()
	ferror      func(error)
	on          map[string]reflect.Value
	ackfunc     map[uint]reflect.Value
	id          uint
}

//OnConnect register a func which is called when connected..
func (c *Client) OnConnect(f func()) {
	c.fconnect = f
}

//OnDisconnect register a func which is called when disconnected.
func (c *Client) OnDisconnect(f func()) {
	c.fdisconnect = f
}

//OnError register a func when error
func (c *Client) OnError(f func(error)) {
	c.ferror = f
}

//On register a func which for event "name".
func (c *Client) On(name string, f interface{}) error {
	c.RLock()
	defer c.RUnlock()
	v := reflect.ValueOf(f)
	t := v.Type()
	if v.Kind() != reflect.Func {
		return errors.New("the last arg must be func")
	}
	if t.NumIn() != 1 {
		return errors.New("the func must have one arg")
	}
	if t.NumOut() != 1 {
		return errors.New("the func must return one result")
	}
	c.on[name] = v
	return nil
}

//Emit emits "name" event with  dat.
func (c *Client) Emit(name string, dat interface{}, f interface{}) error {
	c.Lock()
	defer c.Unlock()
	p := &packetWrite{
		Type: headerEvent,
		ID:   c.id,
		Param: &eventParamWrite{
			Name:  name,
			Param: dat,
		},
	}
	c.send <- p
	if f == nil {
		return nil
	}
	v := reflect.ValueOf(f)
	t := v.Type()
	if v.Kind() != reflect.Func {
		return errors.New("the last arg must be func")
	}
	if t.NumIn() != 1 {
		return errors.New("the func must have one arg")
	}
	if t.NumOut() != 0 {
		return errors.New("the func must  return nothing")
	}
	c.ackfunc[c.id] = v
	c.id++
	return nil
}

func (c *Client) error(err error) {
	c.Lock()
	defer c.Unlock()
	c.send <- &packetWrite{
		Type:  headerError,
		ID:    c.id,
		Param: []byte(err.Error()),
	}
	c.id++
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() error {
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return errors.WithStack(err)
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		var p packetRead
		if err := c.conn.ReadJSON(&p); err != nil {
			return errors.WithStack(err)
		}
		switch p.Type {
		case headerConnect:
			if c.fconnect != nil {
				c.fconnect()
			}
		case headerEvent:
			var param eventParamRead
			if err := json.Unmarshal(p.Param, &param); err != nil {
				c.error(err)
				log.Println(err)
				continue
			}
			f, ok := c.on[param.Name]
			if !ok {
				err := errors.New(param.Name + "not registered")
				c.error(err)
				log.Println(err)
				continue
			}
			t := f.Type()
			obj := reflect.New(t.In(0)).Interface()
			if err := json.Unmarshal(param.Param, obj); err != nil {
				c.error(err)
				log.Println(err)
				continue
			}
			val := reflect.Indirect(reflect.ValueOf(obj))
			result := f.Call([]reflect.Value{val})
			c.send <- &packetWrite{
				Type:  headerACK,
				ID:    p.ID,
				Param: result[0].Interface(),
			}
		case headerACK:
			f, ok := c.ackfunc[p.ID]
			if !ok {
				// log.Println(p.ID, "no callback")
				continue
			}
			t := f.Type()
			obj := reflect.New(t.In(0)).Interface()
			if err := json.Unmarshal(p.Param, obj); err != nil {
				c.error(err)
				log.Println(err)
				continue
			}
			val := reflect.Indirect(reflect.ValueOf(obj))
			f.Call([]reflect.Value{val})
		case headerError:
			if c.ferror != nil {
				c.ferror(errors.New(string(p.Param)))
			}
		default:
			err := errors.Errorf("%d unknown type", p.Type)
			c.error(err)
			log.Println(err)
			continue
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(ctx context.Context) error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx2.Done():
			return nil
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return errors.WithStack(err)
			}
			if !ok {
				return errors.WithStack(c.conn.WriteMessage(websocket.CloseMessage, []byte{}))
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return errors.WithStack(err)
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return errors.WithStack(err)
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return errors.WithStack(err)
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(client *Client, w http.ResponseWriter, r *http.Request) {
	client.Lock()
	defer client.Unlock()
	if client.conn != nil {
		log.Println("only one person can connect")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client.conn = conn

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	wg.Add(2)
	go func() {
		if err := client.writePump(ctx); err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	go func() {
		if err := client.readPump(); err != nil {
			log.Println(err)
		}
		wg.Done()
		cancel()
	}()
	go func() {
		wg.Wait()
		if err := client.conn.Close(); err != nil {
			log.Println(err)
		}
		log.Println("closed")
		client.conn = nil
		time.Sleep(3 * time.Second)
		if client.conn == nil && client.ferror != nil {
			client.ferror(errors.New("connection closed"))
			log.Println("end of pump")
		}
	}()
	client.send <- &packetWrite{
		Type: headerConnect,
		ID:   client.id,
	}
	client.id++
}

//Close close connection.
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		log.Println(err)
	}
}

func newClient() *Client {
	return &Client{
		send:    make(chan *packetWrite, 8),
		on:      make(map[string]reflect.Value),
		ackfunc: make(map[uint]reflect.Value),
	}
}
