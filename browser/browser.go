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

package browser

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

var optionAppStyleChrome = []string{"--disable-extension", "--new-window", "--app=%s"}

type cmdarg struct {
	cmd []string
	arg []string
}

//replaceURL replaces %s in orig to url.
func (c *cmdarg) replaceURL(url string) {
	for i, o := range c.arg {
		if strings.Contains(o, "%s") {
			c.arg[i] = fmt.Sprintf(o, url)
		}
	}
}

//browsers returns the browser paths.
func browsers() []*cmdarg {
	cmds := []*cmdarg{chromePaths()}
	for _, c := range cmds {
		c.arg = append(c.arg, optionAppStyleChrome...)
	}
	cmds = append(cmds, defaultPaths())
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, &cmdarg{
			cmd: []string{exe},
			arg: []string{"%s"},
		})
	}
	return cmds
}

//ErrNoBrowser is an error that no browsers are found.
var ErrNoBrowser = errors.New("no browsers are found")

//Start  starts browsers and opens URL p.
func Start(p string) error {
	cmds := browsers()
	for _, cmd := range cmds {
		cmd.replaceURL(p)
	}
	for _, c := range cmds {
		for _, cmd := range c.cmd {
			viewer := exec.Command(cmd, c.arg...)
			//viewer.Stderr = os.Stderr
			if err := viewer.Start(); err != nil {
				continue
			}
			return nil
		}
	}
	return ErrNoBrowser
}
