
interface Message {
    x: number,
    y: number,
    z: number    
}

interface Message extends Array<Message>{}

class Connection {

    private socket: WebSocket;
    private target: any;

    constructor(target: any) {
        this.open = this.open.bind(this);
        this.message = this.message.bind(this);

        this.target = target;

        this.socket = new WebSocket('ws://localhost:1337')
        this.socket.onopen = this.open;
        this.socket.onmessage = this.message;
    }

    private open() {
        console.log('connection ok');
    }

    private message(payload: MessageEvent) {
        let data: Message[];

        //console.log(payload.data);

        try {
            data = JSON.parse(payload.data)
        } catch(err) {
            console.error(err);
            return;
        }
        
        this.target(data);

        //console.log(data);
    }

}

export default Connection;