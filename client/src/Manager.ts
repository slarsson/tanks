import Game from './Game';
import Connection from './Connection';
import Graphics from './Graphics';

class Manager {

    public static readonly BROADCAST_RATE: number = 50;

    private wasm: any;
    private game: Game;
    private connection: Connection;
    private graphics: Graphics;
    private playerNames: Map<number, string>;

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
        this.playerNames = new Map();
    }

    private init() {
        console.log('MANAGER: connection ok, start game');

        // init send-loop
        this.pollState();

        // add listeners
        window.addEventListener('keydown', this.registerKey);
        window.addEventListener('keyup', this.unregisterKey);

        // connect to game session
        let arr = new Uint8Array(1)
        arr[0] = 0;
        this.connection.send(arr.buffer);

        // input player name
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

        switch (mt) {
            case 9: {
                let id = (new Uint32Array(buffer))[0];
                console.log('remove ID:', id);
                this.wasm.removePlayer(id);
                this.game.removePlayer(id);
                this.playerNames.delete(id);
                this.graphics.setConnectedPlayers(Array.from(this.playerNames.values()));
            }break;

            case 10: {
                let id = (new Uint32Array(buffer))[0];
                this.wasm.setSelf(id);
                this.game.setSelf(id);
            } break; 

            case 99: {
                let data = new Uint32Array(buffer);
                let k1 = this.playerNames.get(data[0]);
                let k2 = this.playerNames.get(data[1]);
                if (k1 == undefined || k2 == undefined) {
                    console.warn('MANAGER: can not show kill log :(');
                } else {
                    this.graphics.addKillMessage(k1, k2);
                }
            } break;
            
            case 98: {
                let id = new Uint32Array(buffer.slice(0, 4))[0];
                let name = new TextDecoder('utf-8').decode(new Uint8Array(buffer.slice(4)));
                this.playerNames.set(id, name);
                
                if (id == this.game.getSelf()) {                
                    this.graphics.removeNameInput();
                    this.game.camera.setFlyTo(0, -40, 30, 500);
                    this.game.camera.setMode(2);
                } 
    
                this.graphics.addMessage(`New player added: {id: ${id}, name: ${name}}`, 2000);
                this.graphics.setConnectedPlayers(Array.from(this.playerNames.values()));
            } break;

            case 33: {
                console.log('GAME: SERVER ERROR ??');
                this.graphics.addMessage('SERVER ERROR', null, Graphics.ERROR);
            } break;

            case 44: {
                let data = new Float32Array(buffer);
                this.game.gameMap.crane.setPosition(data[0], data[1], 0);
                this.wasm.updateCrane(data[0], data[1]);
            } break;

            default: console.log('MANAGER: unknown message', mt);
        }
    }

    private errorHandler(message: string): void {
        this.graphics.addMessage(message, null, Graphics.ERROR);
    }

    private pollState(): void {
        this.connection.send(this.wasm.poll().buffer);
        setTimeout(this.pollState, Manager.BROADCAST_RATE);
    }

    private registerKey(evt: KeyboardEvent): void {
        if (!this.wasm.keypress(evt.key, true)) {
            //console.log(`MANAGER: key not used by wasm (${evt.key})`);
            
            // add key-events here..

            if (evt.key == 'CapsLock') {
                this.graphics.addMessage('CAPS LOCK WARNING :(', 2000, Graphics.WARNING);
                return;
            }
        }
    }

    private unregisterKey(evt: KeyboardEvent): void {
        if (!this.wasm.keypress(evt.key, false)) {

            // unregister key-events here..
        }
    }
}

export default Manager;