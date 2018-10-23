// +build !windows

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

package browser

import "runtime"

//DefaultPath returns paths of default browsers.
func defaultPaths() *cmdarg {
	switch runtime.GOOS {
	case "darwin":
		return &cmdarg{
			cmd: []string{"/usr/bin/open"},
			arg: []string{"%s"},
		}

	default:
		return &cmdarg{
			cmd: []string{"xdg-open"},
			arg: []string{"%s"},
		}
	}
}

//ChromePath returns paths of chrome.
func chromePaths() *cmdarg {
	switch runtime.GOOS {
	case "darwin":
		return &cmdarg{
			cmd: []string{"/usr/bin/open"},
			arg: []string{"-n", "-a", "Google Chrome", "--args", "%s"},
		}

	default:
		return &cmdarg{
			cmd: []string{
				"chrome",
				"google-chrome",
				"chrome-stable",
				"google-chrome-stable",
				"/opt/google/chrome/chrome",
				"/opt/google/chrome/google-chrome",
			},
			arg: []string{"%s"},
		}
	}
}
