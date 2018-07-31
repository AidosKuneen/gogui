[![Build Status](https://travis-ci.org/AidosKuneen/gogui.svg?branch=master)](https://travis-ci.org/AidosKuneen/gogui)
[![GoDoc](https://godoc.org/github.com/AidosKuneen/gogui?status.svg)](https://godoc.org/github.com/AidosKuneen/gogui)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/AidosKuneen/gogui/master/LICENSE)


# gogui

This library is a UI alternative for Go which uses the Chrome browser (or a default browser) which is already installed. 
The application communicates with the browser via [socket.io](https://socket.io/).
 That means you can call funcs in the browser from the server (and vice versa) without worrying about the websocket and JavaScript.

You can write the server-side in Go and client-side in Javascript.

## Requirements

This requires

* git
* go 1.10+
* web browser
	* Newest Chrome browser is recommended.
	* firefox
	* Microsoft Edge

## Platforms

* Linux
* OSX
* Windows

## Installation

    $ go get -u github.com/AidosKuneen/gogui

## Example

See [example directory](https://github.com/AidosKuneen/gogui/example)

# Contribution
Improvements to the codebase and pull requests are encouraged.

