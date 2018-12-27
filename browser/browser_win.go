// +build windows

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
	"log"

	"golang.org/x/sys/windows/registry"
)

//DefaultPath returns paths of default browsers.
func defaultPaths() *cmdarg {
	return &cmdarg{
		cmd: []string{"cmd"},
		arg: []string{"/c", "start", "%s"},
	}
}

//ChromePath returns paths of chrome.
func chromePaths() *cmdarg {
	regpath := `SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, regpath, registry.QUERY_VALUE)
	if err != nil {
		log.Println(regpath, err)
		return nil
	}
	defer k.Close()

	s, _, err := k.GetStringValue("")
	if err != nil {
		log.Println(err)
		return nil
	}
	return &cmdarg{
		cmd: []string{s},
		arg: []string{"%s"},
	}
}
