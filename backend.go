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
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/AidosKuneen/gogui/browser"
	"github.com/pkg/errors"
)

//GUI is chans for notifiing connected and finished.
type GUI struct {
	Finished  chan error
	Connected chan struct{}
	Client    *Client
}

//New returns a GUI struct.
func New() *GUI {
	return &GUI{
		Finished:  make(chan error, 2),
		Connected: make(chan struct{}, 100), //for reload
		Client:    newClient(),
	}
}

//On registers handler for event n.
func (g *GUI) On(n string, f interface{}) {
	if err := g.Client.On(n, f); err != nil {
		panic(err)
	}
}

//Emit emits "name" event with  dat.
func (g *GUI) Emit(name string, dat interface{}, f interface{}) error {
	return g.Client.Emit(name, dat, f)
}

//Start starts GUI bakcend.
//Set dest to react debug server URL for redirecting to it.
//You must setup http.Handle before calling it.
func (g *GUI) Start(dest, path string) error {
	g.Client.OnError(func(e error) {
		log.Println("error", e)
		if len(g.Finished) == 0 {
			g.Finished <- e
		}
	})
	g.Client.OnDisconnect(func() {
		log.Println("closed")
		g.Finished <- nil
	})
	g.Client.OnConnect(func() {
		log.Println("connected")
		g.Connected <- struct{}{}
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(g.Client, w, r)
	})
	if dest != "" {
		http.HandleFunc("/", doProxy(dest))
	}
	pno, err := freePort(54244)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Println("Serving at localhost:", pno, "...")
	go func() {
		defer g.Client.Close()
		if err := http.ListenAndServe(":"+strconv.Itoa(pno), nil); err != nil {
			log.Println(err)
			g.Finished <- err
		}
	}()
	//disable cache
	v := strconv.Itoa(int(time.Now().Unix()))
	return browser.Start("http://localhost:" + strconv.Itoa(pno) + "/" + path + "?v=" + v)
}

func doProxy(dest string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp *http.Response
		client := &http.Client{}
		url := dest + r.RequestURI
		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for name, value := range r.Header {
			req.Header.Set(name, value[0])
		}
		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = r.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for k, v := range resp.Header {
			w.Header().Set(k, v[0])
		}
		w.WriteHeader(resp.StatusCode)
		if _, err = io.Copy(w, resp.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = resp.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func freePort(def int) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:"+strconv.Itoa(def))
	if err != nil {
		return 0, errors.WithStack(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	defer func() {
		if err = l.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err == nil {
		return def, nil
	}

	addr, err = net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, errors.WithStack(err)
	}

	l, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return l.Addr().(*net.TCPAddr).Port, nil
}
