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

enum header {
    connect,
    event,
    ack,
    error,
}

class GUI {
    private conn: WebSocket;
    private fconnect: () => void;
    private fdisconnect: () => void;
    private ferror: (arg: string) => void;
    private fon: { [key: string]: (arg: any) => any } = {};
    private ackfunc: { [key: number]: (arg: any) => void } = {};
    private id = 0;

    public onConnect = (f: () => void) => {
        this.fconnect = f
    }
    public onDisconnect = (f: () => void) => {
        this.fdisconnect = f
    }
    public onError = (f: (arg: string) => void) => {
        this.ferror = f
    }
    public on = (name: string, f: (arg: any) => any) => {
        this.fon[name] = f
    }
    public emit = (name: string, dat: any, f: (arg: any) => void) => {
        const msg = {
            id: this.id,
            param: {
                name,
                param: dat,
            },
            type: header.event,
        }
        this.ackfunc[this.id] = f
        this.id++
        this.conn.send(JSON.stringify(msg))
    }

    public connect = () => {
        const parent = this
        if (!document || !document.location){
            console.log("page not loaded")
            return
        }
        this.conn = new WebSocket("ws://" + document.location.host + "/ws");
        this.conn.onopen = () => {
            const msg = {
                id: parent.id,
                type: header.connect,
            }
            parent.id++
            parent.conn.send(JSON.stringify(msg))
        }
        this.conn.onclose = (ev: CloseEvent) => {
            if (this.onError) {
                this.ferror(ev.reason)
            }
            if (this.onDisconnect) {
                this.fdisconnect()
            }
        }
        this.conn.onmessage = (event: MessageEvent) => {
            const msg = JSON.parse(event.data);
            switch (msg.type) {
                case header.connect:
                    if (parent.fconnect) {
                        parent.fconnect()
                    }
                    break
                case header.event:
                    {
                        const f = parent.fon[msg.param.name]
                        if (f == null) {
                            break;
                        }
                        const result = f(msg.param.param)
                        const resp = {
                            id: parent.id,
                            param: result,
                            type: header.ack,
                        }
                        parent.conn.send(JSON.stringify(resp))
                    }
                    break
                case header.ack:
                    {
                        const f = parent.ackfunc[msg.id]
                        if (f == null) {
                            break;
                        }
                        f(msg.param)
                    }
                    break
                case header.error:
                    if (parent.onError) {
                        parent.ferror(msg.param)
                    }
                    break
            };
        }
    }
}


export default GUI;
