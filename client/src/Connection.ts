
interface Message {
    x: number,
    y: number,
    z: number    
}

interface Message extends Array<Message>{}

class Connection {

    private socket: WebSocket;
    private target: any;

    constructor(target: any, swag: any) {
        //this.open = this.open.bind(this);
        this.message = this.message.bind(this);
        this.close = this.close.bind(this);

        this.target = target;

        this.socket = new WebSocket('ws://157.245.70.83:1337')
        //this.socket = new WebSocket('ws://localhost:1337')
        this.socket.binaryType = 'arraybuffer';
        this.socket.onopen = () => swag();
        this.socket.onmessage = this.message;
        this.socket.onclose = this.close;
    }

    // private open() {
    //     console.log('connection ok');
    // }

    private close() {
        console.log('connection closed');
    }

    private message(payload: MessageEvent) {
        this.target(new Uint32Array(payload.data.slice(0, 4)), payload.data.slice(4));
        return;
        
        //console.log(payload.data.Uint8Array);
        let messageType = new Uint32Array(payload.data.slice(0, 4));
        
        //if (messageType[0] < 100) {
            //console.log('messageType:', messageType);
        //}
        // console.log( new Uint32Array(payload.data));
        // console.log(new Float32Array(buf));
        if (messageType[0] == 0) {
            console.log('fan ta postnord');
            this.target(new Float32Array(payload.data.slice(4)), new Uint32Array());
        } else {
            console.log(payload.data);
            this.target(new Float32Array(), new Uint32Array(payload.data));
        }

        // let x = new Uint32Array(payload.data);

        // for (var i = 0; i < x.byteLength; i++) {
        //     console.log('wtf?');
        // }

        //this.target();
        //this.target(String.fromCharCode.apply(null, x));
        // let data: Message[];

        // try {
        //     data = JSON.parse(payload.data)
        // } catch(err) {
        //     console.error(err);
        //     return;
        // }
        
        // this.target(data);
    }

    send(message: any) {
        this.socket.send(message);
        //console.log('send message:', message);
    }

}

export default Connection;