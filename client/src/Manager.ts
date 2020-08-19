import Game from './Game';
import Connection from './Connection2';
import Graphics from './Graphics';

class Manager {

    public static readonly BROADCAST_RATE: number = 50;

    private wasm: any; // TODO: fix 'any'
    private game: Game;
    private connection: Connection;
    private graphics: Graphics;

    constructor(wasm: any) {
        this.messageHandler = this.messageHandler.bind(this);
        this.errorHandler = this.errorHandler.bind(this);
        this.init = this.init.bind(this);
        this.pollState = this.pollState.bind(this);
        this.registerKey = this.registerKey.bind(this);
        this.unregisterKey = this.unregisterKey.bind(this);
        
        this.wasm = wasm;
        this.graphics = new Graphics();
        this.game = new Game(this.wasm, this.graphics);
        this.connection = new Connection('ws://localhost:1337', this.messageHandler, this.init, this.errorHandler);
        

        console.log('my manager..');
    }

    private init() {
        console.log('MANAGER: connection ok, start game');

        //this.graphics.addMessage('Connection OK', 2000);
        this.pollState();

        window.addEventListener('keydown', this.registerKey);
        window.addEventListener('keyup', this.unregisterKey);
        //this.graphics.newNameInput((x: string) => console.log(x));

        // MOVE THIS SOMEWHERE ELSE ..
        let arr = new Uint8Array(1)
        arr[0] = 0;
        this.connection.send(arr.buffer);

        // SWAG
        this.graphics.newNameInput((s: string) => {
            let arr = new Uint8Array(1 + s.length);
            arr[0] = 99;

            for (let i = 0; i < s.length; i++) {
                arr[i + 1] = s[i].charCodeAt(0); 
            }

            this.connection.send(arr.buffer);
        });
    }

    private messageHandler(mt: number, buffer: ArrayBuffer): void {
        if (mt == 0) {
            let state = new Float32Array(buffer);   
            this.wasm.update(...state);
            this.game.update(state);
            return;
        }

        console.log('MANAGER: new message:', mt);

        switch (mt) {
            case 9: 
                {
                    let id = (new Uint32Array(buffer))[0];
                    console.log('remove ID:', id);
                    this.wasm.removePlayer(id);
                    this.game.removePlayer(id);
                }
                break;

            case 10:
                let id = (new Uint32Array(buffer))[0];
                this.wasm.setSelf(id);
                this.game.setSelf(id);
                break;

            case 99:
                {
                    let data = new Uint32Array(buffer);
                    console.log(data);
                    console.log('show kill log ?');
        
                    // let k1 = "wtfplayer" + Math.random();
                    // let k2 = "asdf" + Math.random();
                    let k1 = 'a';
                    let k2 = 'b';
        
                    if (k1 != undefined && k2 != undefined) {
                        this.graphics.addKillMessage(k1, k2);
                    }
                }
                break;
            
            case 98:
                {
                    let id = new Uint32Array(buffer.slice(0, 4));
                    let name = new TextDecoder('utf-8').decode(new Uint8Array(buffer.slice(4)));
                    //this.players.set(id[0], name);
        
                    console.log('self:', this.game.getSelf());
                    if (id[0] == this.game.getSelf()) {                
                        this.graphics.removeNameInput();
                        
                        this.game.camera.setFlyTo(0, -40, 30, 500);
                        this.game.camera.setMode(2);
                    } 
        
                    this.graphics.addMessage(`New player added: {id: ${id}, name: ${name}}`, 2000);
        
                    // console.log('id:', new Uint32Array(buffer.slice(0, 4)));
                    // console.log(new TextDecoder('utf-8').decode(new Uint8Array(buffer.slice(4))));

                }
                break;

            case 33:
                console.log('GAME: SERVER ERROR ??');
                this.graphics.addMessage('SERVER ERROR', null, Graphics.ERROR);
                break;
        }
    }

    private errorHandler(): void {
        console.log('my error...');
    }

    private pollState(): void {
        this.connection.send(this.wasm.poll().buffer);
        setTimeout(this.pollState, Manager.BROADCAST_RATE);
    }

    private registerKey(evt: KeyboardEvent): void {
        if (!this.wasm.keypress(evt.key, true)) {
            console.log(`MANAGER: key not used by wasm (${evt.key})`);
            if (evt.key == 'CapsLock') {
                this.graphics.addMessage('CAPS LOCK WARNING :(', 2000, Graphics.WARNING);
                return;
            }

            if (evt.key == 'z') {
                this.game.camera.setMode(2);
            }
        }
    }

    private unregisterKey(evt: KeyboardEvent): void {
        this.wasm.keypress(evt.key, false);
    }
}

export default Manager;