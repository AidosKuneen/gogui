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

package main

import (
	"log"
	"os"
	"time"

	"github.com/AidosKuneen/gogui"
	"gopkg.in/googollee/go-socket.io.v1"
)

func main() {
	funcs := map[string]interface{}{
		"msg": func(s socketio.Conn, msg string) string {
			log.Println("receive message", msg)
			s.Emit("reply", "emit from server "+msg, func(dat string) {
				log.Println("emit from server ", dat)
			})
			return "ack " + msg
		},
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	dest := ""
	if len(os.Args) > 1 {
		dest = os.Args[1]
	}
	gui, err := gogui.Start(funcs, dest)
	if err != nil {
		log.Fatal(err)
	}
	var con socketio.Conn
	select {
	case <-time.After(60 * time.Second):
		log.Println("failed to initialize")
		return
	case con = <-gui.Connected:
	}
	con.Emit("reply", "after 3 secs", func(dat string) {
		log.Println("emit from server ", dat)
	})
	if err := <-gui.Finished; err != nil {
		log.Println(err)
	}
}
