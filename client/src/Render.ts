import * as THREE from 'three';
import Tank from './Tank';
//import Vehicle from './Vehicle';
import Keypress from './Keypress';
import Connection from './Connection';
import GameMap from './Map'
import Graphics from './Graphics';

import Projectile from './Particle';


const testName = (conn: Connection) => {
    console.log('meh...', conn);

    let root = document.getElementById('test');
    let container = document.createElement('div');
    container.classList.add('input-test');

    let form = document.createElement('form');

    let input = document.createElement('input');
        input.setAttribute('type', 'text');
        
    container.appendChild(form);
    form.appendChild(input);
    //root?.appendChild(container);


    form.addEventListener('submit', (e: any) => {
        e.preventDefault();
        console.log(input.value);
        
        let s = input.value;
            
        let arr = new Uint8Array(1 + s.length);
        arr[0] = 99;

        for (let i = 0; i < s.length; i++) {
            //onsole.log(s[i]);
            arr[i + 1] = s[i].charCodeAt(0); 
        }

        conn.send(arr.buffer);
    });
};



class Render {

    public readonly BROADCAST_RATE = 50;

    private wasm: any;
    private conn: Connection;
    private scene: THREE.Scene;
    private camera: THREE.PerspectiveCamera;
    private renderer: THREE.WebGLRenderer;
    private timestamp: number;
    
    private vehicles: Map<number, Tank>;
    private self: number = -1;
    private shoot: Map<number, Projectile>;

    private SWAG: boolean;

    private gg: Graphics;


    private players: Map<number, string>;

    constructor(wasm: any) {
        this.SWAG = false;
        
        // wasm => object containing go functions
        this.wasm = wasm;
        
        // method binds
        this.animate = this.animate.bind(this);
        this.broadcast = this.broadcast.bind(this);
        this.serverMessage = this.serverMessage.bind(this);
        this.registerKey = this.registerKey.bind(this);
        this.unregisterKey = this.unregisterKey.bind(this);

        // setup scene + camera + renderer
        this.scene = new THREE.Scene();

        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth/window.innerHeight, 1, 2000);
        //this.camera.rotation.x = Math.PI / 4;
        this.camera.position.y = -140;
        this.camera.position.z = 30;
        this.camera.lookAt(0, 0, 0);
        
        this.renderer = new THREE.WebGLRenderer({antialias: true, alpha: false});
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        this.renderer.setClearColor(0x000000, 0);
        
        // setup objects + testing..
        this.vehicles = new Map();
        this.shoot = new Map();
        this.players = new Map();
        new GameMap(this.scene);
        
        this.gg = new Graphics();
        this.gg.newInfoBox('Postnord ska och borde läggas ner, typ så..', null, Graphics.INFO);
        this.gg.newInfoBox('lägg ner postnord igen..', null, Graphics.ERROR);
        this.gg.newInfoBox('Warning my warning.', null, Graphics.WARNING);

        this.gg.newKillMessage('jocke', 'jonna');
        //this.shoot = new Particle(this.scene);


        // init websocket connection
        this.conn = new Connection(
            this.serverMessage,
            () => {
                // test..
                let arr = new Uint8Array(1)
                arr[0] = 0;
                this.conn.send(arr.buffer);
                this.broadcast();
            }
        );

        // setup eventlisteners
        window.addEventListener('keydown', this.registerKey);
        window.addEventListener('keyup', this.unregisterKey);

        // init, add to DOM and init renderloop
        document.getElementById('root')?.appendChild(this.renderer.domElement);
        this.timestamp = performance.now();
        this.animate();

        // test:


        testName(this.conn);
        // setTimeout(() => {
        //     let s = "##läggnerpostnord"
            
        //     let arr = new Uint8Array(1 + s.length);
        //     arr[0] = 99;

        //     for (let i = 0; i < s.length; i++) {
        //         console.log(s[i]);
        //         arr[i + 1] = s[i].charCodeAt(0); 
        //     }

        //     this.conn.send(arr.buffer);
        // }, 1000);
    }

    private serverMessage(mt: number, buffer: ArrayBuffer): void {
        //console.log(mt);
        if (mt == 0) {
            let state = new Float32Array(buffer);            
            for (let i = 0; i < state.length; i += 12) {
                if (!this.vehicles.has(state[i])) {
                    this.vehicles.set(state[i], new Tank(this.scene));
                } else if (state[i+9] == 1) {
                    console.log('ADD ANOTHER PROJECTILE FFS');
                    //this.wasm.addProjectile(state[i]); // float comparsion missmatch!?
                    // float comparision error? wtf?
                    // console.log('SHOOTING FFS');
                    // const v = this.vehicles.get(this.self);
                    // if (v != undefined) {
                    //     this.shoot.add2(v.getGunRotation());
                    // }
                } else if (state[i+11] == 0) {
                    console.log('player is dead');
                    this.vehicles.get(state[i])?.kill();
                } else {
                    this.vehicles.get(state[i])?.respawn();
                }  
            }
            this.wasm.update(...state, this.SWAG);
        }
        
        if (mt == 10) {
            let test = new Uint32Array(buffer);
            console.log('my ID:', test);
            this.self = test[0];
            this.wasm.setSelf(this.self);
            return;
        } 
        
        if (mt == 9) {
            let test = new Uint32Array(buffer);
            console.log('remove ID:', test);
            this.wasm.removePlayer(test[0]);
            //this.scene.remove(this.vehicles.get(test[1]))
            this.vehicles.get(test[0])?.dispose();
            this.vehicles.delete(test[0]);
            return;
        }
        
        if (mt == 99) {
            let data = new Uint32Array(buffer);
            console.log(data);
            console.log('show kill log ?');

            let k1 = "wtfplayer" + Math.random();
            let k2 = "asdf" + Math.random();
            // let k1 = this.players.get(data[0]);
            // let k2 = this.players.get(data[1]);

            if (k1 != undefined && k2 != undefined) {
                this.gg.newKillMessage(k1, k2);
            }

            // let div = document.createElement('div');
            // div.innerHTML = '<span>'+this.players.get(data[0])+'</span> KILLED <span>'+this.players.get(data[1])+'</span>';
            // document.getElementById('kills')?.appendChild(div);

            // setTimeout(() => {
            //     div.remove()
            // }, 1000);
        }

        if (mt == 98) {
            let id = new Uint32Array(buffer.slice(0, 4));
            let name = new TextDecoder('utf-8').decode(new Uint8Array(buffer.slice(4)));
            this.players.set(id[0], name);

            if (id[0] == this.self) {
                let div = document.getElementById('test');
                if (div != null) {
                    div.innerHTML = '';
                }
            } 

            // console.log('id:', new Uint32Array(buffer.slice(0, 4)));
            // console.log(new TextDecoder('utf-8').decode(new Uint8Array(buffer.slice(4))));
        }

        if (mt == 33) {
            console.log('SERVER ERROR WTF');
        }
    }

    // private serverMessage(state: Float32Array, test: Uint8Array): void {
    //     // //console.log(state);
    //     // if (test[0] == 10) {
    //     //     console.log('my ID:', test);
    //     //     this.self = test[1];
    //     //     this.wasm.setSelf(this.self);
    //     //     //this.vehicles.get(test[1])?.setColor(0xf699cd);

    //     //     return;
    //     // } 
        
    //     // if (test[0] == 9) {
    //     //     console.log('remove ID:', test);
    //     //     //this.scene.remove(this.vehicles.get(test[1]))
    //     //     this.vehicles.get(test[1])?.dispose();
    //     //     this.vehicles.delete(test[1]);
    //     //     return;
    //     // } 
        
    //     // this.wasm.update(...state);
    
    //     // for (let i = 0; i < state.length; i += 10) {
    //     //     if (!this.vehicles.has(state[i])) {
    //     //         this.vehicles.set(state[i], new Vehicle(this.scene));
    //     //     }        
    //     // }
    // }

    private broadcast(): void {
        //this.wasm.poll()
        //console.log(this.wasm.poll());
        this.conn.send(this.wasm.poll().buffer);
        setTimeout(this.broadcast, this.BROADCAST_RATE);
    }

    private registerKey(evt: KeyboardEvent): void {
        this.wasm.keypress(evt.key, true);

        if (evt.key == 'c') {
            //this.vehicles.get(this.self)?.setColor(0xfc99cd);
            this.SWAG = !this.SWAG;
        }

        // if (evt.key == ' ') {
        //     if (this.self != -1) {
        //         const v = this.vehicles.get(this.self);
        //         if (v != undefined) {
        //             //this.shoot.add(v.getGunOrigin());
        //             this.shoot.add2(v.getGunRotation());
        //         }
        //     }
        // }
    }

    private unregisterKey(evt: KeyboardEvent): void {
        this.wasm.keypress(evt.key, false);
    }

    public animate(): void {
        requestAnimationFrame(this.animate);
        let now = performance.now();
        let dt = now - this.timestamp;

        //this.shoot.update(dt);
        
        //this.wasm.guessPosition(dt);


        const it = this.vehicles[Symbol.iterator]();
        for (let item of it) {
            const pos = this.wasm.getPos(item[0], dt)
            // //console.log(pos);
            // if (pos.length == 0) {
            //     item[1].dispose()
            //     continue
            // }

            item[1].setPosition(pos[0], pos[1], pos[2]);
            item[1].setRotation(pos[3]);
            item[1].setTurretRotation(pos[4]);
            
            if (item[0] == this.self) {
                this.camera.position.x = pos[0];
                this.camera.position.y = pos[1] - 40;
            }
        }

        let items = this.wasm.updateProjectiles(dt);
        for (let i = 0; i < items.length; i += 5) {
            if (items[i+4] == 0) {
                //console.log('remove:', items[i]);
                this.shoot.get(items[i])?.dispose();
                this.shoot.delete(items[i]);
                continue
            } 
            
            if (!this.shoot.has(items[i])) {
                this.shoot.set(items[i], new Projectile(this.scene));
            }
            this.shoot.get(items[i])?.set(items[i+1], items[i+2], items[i+3]);
        }
        //console.log(this.wasm.updateProjectiles(dt));

        this.renderer.render(this.scene, this.camera);    
        this.timestamp = now;
    }

}

export default Render;