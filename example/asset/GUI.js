"use strict";
// Copyright (c) 2018 Aidos Developer
exports.__esModule = true;
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
var header;
(function (header) {
    header[header["connect"] = 0] = "connect";
    header[header["event"] = 1] = "event";
    header[header["ack"] = 2] = "ack";
    header[header["error"] = 3] = "error";
})(header || (header = {}));
var GUI = /** @class */ (function () {
    function GUI() {
        var _this = this;
        this.fon = {};
        this.ackfunc = {};
        this.id = 0;
        this.onConnect = function (f) {
            _this.fconnect = f;
        };
        this.onDisconnect = function (f) {
            _this.fdisconnect = f;
        };
        this.onError = function (f) {
            _this.ferror = f;
        };
        this.on = function (name, f) {
            _this.fon[name] = f;
        };
        this.emit = function (name, dat, f) {
            var msg = {
                id: _this.id,
                param: {
                    name: name,
                    param: dat
                },
                type: header.event
            };
            _this.ackfunc[_this.id] = f;
            _this.id++;
            _this.conn.send(JSON.stringify(msg));
        };
        this.connect = function () {
            var parent = _this;
            if (!document || !document.location) {
                console.log("page not loaded");
                return;
            }
            _this.conn = new WebSocket("ws://" + document.location.host + "/ws");
            _this.conn.onopen = function () {
                var msg = {
                    id: parent.id,
                    type: header.connect
                };
                parent.id++;
                parent.conn.send(JSON.stringify(msg));
            };
            _this.conn.onclose = function (ev) {
                if (_this.onError) {
                    _this.ferror(ev.reason);
                }
                if (_this.onDisconnect) {
                    _this.fdisconnect();
                }
            };
            _this.conn.onmessage = function (event) {
                var msg = JSON.parse(event.data);
                switch (msg.type) {
                    case header.connect:
                        if (parent.fconnect) {
                            parent.fconnect();
                        }
                        break;
                    case header.event:
                        {
                            var f = parent.fon[msg.param.name];
                            if (f == null) {
                                break;
                            }
                            var result = f(msg.param.param);
                            var resp = {
                                id: parent.id,
                                param: result,
                                type: header.ack
                            };
                            parent.conn.send(JSON.stringify(resp));
                        }
                        break;
                    case header.ack:
                        {
                            var f = parent.ackfunc[msg.id];
                            if (f == null) {
                                break;
                            }
                            f(msg.param);
                        }
                        break;
                    case header.error:
                        if (parent.onError) {
                            parent.ferror(msg.param);
                        }
                        break;
                }
                ;
            };
        };
    }
    return GUI;
}());
exports["default"] = GUI;
