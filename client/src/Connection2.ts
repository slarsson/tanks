class Connection2 {

    private socket: WebSocket;
    private target: any;

    constructor(host: string, target: any, onsuccess: any, onerror: any) {
        this.message = this.message.bind(this);
        this.close = this.close.bind(this);

        this.target = target;

        this.socket = new WebSocket(host);

        //this.socket = new WebSocket('ws://samme.international:1337')
        //this.socket = new WebSocket('ws://localhost:1337')


        this.socket.binaryType = 'arraybuffer';
        this.socket.onopen = onsuccess;
        this.socket.onmessage = this.message;
        this.socket.onclose = this.close;
        this.socket.onerror = onerror;
    }

    // private open() {
    //     console.log('connection ok');
    // }

    private close() {
        console.log('connection closed');
    }

    private message(payload: MessageEvent) {
        this.target((new Uint32Array(payload.data.slice(0, 4))[0]), payload.data.slice(4));
    }

    send(message: any) {
        this.socket.send(message);
    }
}

export default Connection2;