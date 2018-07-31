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

	"gopkg.in/googollee/go-socket.io.v1"
)

//GUI is chans for notifiing connected and finished.
type GUI struct {
	Connected chan socketio.Conn
	Finished  chan error
}

//Start starts GUI bakcend.
//Set dest to react debug server URL for redirecting to it.
func Start(funcs map[string]interface{}, dest string) (*GUI, error) {
	gui := &GUI{
		Finished:  make(chan error),
		Connected: make(chan socketio.Conn),
	}
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	for n, f := range funcs {
		server.OnEvent("/", n, f)
	}
	server.OnEvent("/", "initialize_gogui", func(s socketio.Conn) {
		gui.Connected <- s
	})
	server.OnError("/", func(e error) {
		log.Println("error", e)
		gui.Finished <- e
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		gui.Finished <- nil
	})
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected:", s.ID())
		return nil
	})
	go server.Serve()

	http.Handle("/socket.io/", server)
	if dest == "" {
		http.Handle("/", http.FileServer(http.Dir("./asset")))
	} else {
		http.HandleFunc("/", doProxy(dest))
	}

	log.Println("Serving at localhost:5000...")
	go func() {
		defer server.Close()
		if err := http.ListenAndServe(":5000", nil); err != nil {
			log.Println(err)
			gui.Finished <- err
		}
	}()
	return gui, StartBrowser("http://localhost:5000")
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
