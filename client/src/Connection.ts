
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

        this.socket = new WebSocket('ws://localhost:1337')
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
        let data: Message[];

        try {
            data = JSON.parse(payload.data)
        } catch(err) {
            console.error(err);
            return;
        }
        
        this.target(data);
    }

    send(message: any) {
        this.socket.send(message);
        //console.log('send message:', message);
    }

}

export default Connection;