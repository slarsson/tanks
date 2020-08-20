class Connection {

    private socket: WebSocket;
    private target: any;

    constructor(host: string, target: any, onsuccess: any, onerror: any) {
        this.message = this.message.bind(this);
        this.target = target;

        this.socket = new WebSocket(host); 
        this.socket.binaryType = 'arraybuffer';
        this.socket.onopen = onsuccess;
        this.socket.onmessage = this.message;
        this.socket.onclose = () => onerror('CONNECTION CLOSED');
        this.socket.onerror = () => onerror('CONNECTION FAILED');
    }

    private message(payload: MessageEvent) {
        this.target((new Uint32Array(payload.data.slice(0, 4))[0]), payload.data.slice(4));
    }

    send(message: any) {
        this.socket.send(message);
    }
}

export default Connection;