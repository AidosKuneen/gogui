![Golangci lint](https://github.com/AidosKuneen/gogui/actions/workflows/golangci-lint.yml/badge.svg)
![Tests](https://github.com/AidosKuneen/gogui/actions/workflows/tests.yml/badge.svg)
[![GoDoc](https://godoc.org/github.com/AidosKuneen/gogui?status.svg)](https://godoc.org/github.com/AidosKuneen/gogui)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/AidosKuneen/gogui/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/AidosKuneen/gogui)](https://goreportcard.com/report/github.com/AidosKuneen/gogui)

# gogui

This library is a GUI alternative for Go which uses the Chrome browser (or a default browser). 
The application communicates with the browser via WebSocket message using simple JSON.

You can write  server-side code in Go and client-side code in Javascript.

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

See [example directory](https://github.com/AidosKuneen/gogui/tree/master/example)

## Contribution
Improvements to the codebase and pull requests are encouraged.


## Dependencies and Licenses

```
github.com/AidosKuneen/gogui  MIT License
github.com/gorilla/websocket  BSD 2-clause "Simplified" License 
Golang Standard Library       BSD 3-clause License
```
