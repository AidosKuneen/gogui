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
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	dest := ""
	if len(os.Args) > 1 {
		dest = os.Args[1]
	}
	gui := gogui.New()

	gui.On("msg", func(msg string) string {
		log.Println("receive message", msg)
		err2 := gui.Emit("reply", "emit from server "+msg, func(dat string) {
			log.Println("emit from server ", dat)
		})
		if err2 != nil {
			log.Fatal(err2)
		}
		r := "ack " + msg
		return r
	})

	if err := gui.Start(dest); err != nil {
		log.Fatal(err)
	}
	select {
	case <-time.After(60 * time.Second):
		log.Println("failed to initialize")
		return
	case <-gui.Connected:
	}
	msg := "after 3 secs"
	err := gui.Emit("reply", msg, func(dat string) {
		log.Println("emit from server ", dat)
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := <-gui.Finished; err != nil {
		log.Println(err)
	}
}
