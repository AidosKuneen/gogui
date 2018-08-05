// Copyright (c) 2018 Aidos Developer

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
	"net/http"

	"github.com/AidosKuneen/gogui/browser"
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
		Finished:  make(chan error),
		Connected: make(chan struct{}),
		Client:    newClient(),
	}
}

//On registers handler for event n.
func (g *GUI) On(n string, f interface{}) error {
	return g.Client.On(n, f)
}

//Emit emits "name" event with  dat.
func (g *GUI) Emit(name string, dat interface{}, f interface{}) error {
	return g.Client.Emit(name, dat, f)
}

//Start starts GUI bakcend.
//Set dest to react debug server URL for redirecting to it.
func (g *GUI) Start(dest string) error {
	g.Client.OnError(func(e error) {
		log.Println("error", e)
		g.Finished <- e
	})
	g.Client.OnDisconnect(func() {
		g.Finished <- nil
	})
	g.Client.OnConnect(func() {
		log.Println("connected")
		g.Connected <- struct{}{}
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(g.Client, w, r)
	})
	if dest == "" {
		http.Handle("/", http.FileServer(http.Dir("./asset")))
	} else {
		http.HandleFunc("/", doProxy(dest))
	}

	log.Println("Serving at localhost:5000...")
	go func() {
		defer g.Client.Close()
		if err := http.ListenAndServe(":5000", nil); err != nil {
			log.Println(err)
			g.Finished <- err
		}
	}()
	return browser.Start("http://localhost:5000")
}

func doProxy(dest string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp *http.Response
		client := &http.Client{}
		url := dest + r.RequestURI
		req, err := http.NewRequest(r.Method, url, r.Body)
		for name, value := range r.Header {
			req.Header.Set(name, value[0])
		}
		resp, err = client.Do(req)
		r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for k, v := range resp.Header {
			w.Header().Set(k, v[0])
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		resp.Body.Close()
	}
}