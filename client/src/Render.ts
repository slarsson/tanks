import * as THREE from 'three';
import Vehicle from './Vehicle';
import Keypress from './Keypress';
import Connection from './Connection';

import Projectile from './Particle';

import { helper, obstacleTest } from './dev';

class Render {

    public readonly BROADCAST_RATE = 50;

    private wasm: any;
    private conn: Connection;
    private scene: THREE.Scene;
    private camera: THREE.PerspectiveCamera;
    private renderer: THREE.WebGLRenderer;
    private timestamp: number;
    
    private vehicles: Map<number, Vehicle>;
    private self: number = -1;
    private shoot: Map<number, Projectile>;

    private SWAG: boolean;

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
        this.camera.rotation.x = Math.PI / 4;
        this.camera.position.y = -10;
        this.camera.position.z = 30;
        
        this.renderer = new THREE.WebGLRenderer({antialias: true});
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        
        // setup objects + testing..
        this.vehicles = new Map();
        this.shoot = new Map();
        //this.shoot = new Particle(this.scene);
        helper(this.scene);
        obstacleTest(this.scene);

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
    }

    private serverMessage(mt: number, buffer: ArrayBuffer): void {
        if (mt == 0) {
            let state = new Float32Array(buffer);            
            for (let i = 0; i < state.length; i += 11) {
                if (!this.vehicles.has(state[i])) {
                    this.vehicles.set(state[i], new Vehicle(this.scene));
                } else if (state[i+9] == 1) {
                    console.log('ADD ANOTHER PROJECTILE FFS');
                    this.wasm.addProjectile(state[i]); // float comparsion missmatch!?
                    // float comparision error? wtf?
                    // console.log('SHOOTING FFS');
                    // const v = this.vehicles.get(this.self);
                    // if (v != undefined) {
                    //     this.shoot.add2(v.getGunRotation());
                    // }
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
            //this.scene.remove(this.vehicles.get(test[1]))
            this.vehicles.get(test[1])?.dispose();
            this.vehicles.delete(test[1]);
            return;
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
        this.conn.send(this.wasm.poll().buffer);
        setTimeout(this.broadcast, this.BROADCAST_RATE);
    }

    private registerKey(evt: KeyboardEvent): void {
        this.wasm.keypress(evt.key, true);

        if (evt.key == 'c') {
            this.vehicles.get(this.self)?.setColor(0xfc99cd);
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
        
        //this.wasm.local(this.self, dt);
        this.wasm.guessPosition(dt);


        const it = this.vehicles[Symbol.iterator]();
        for (let item of it) {
            const pos = this.wasm.getPos(item[0], dt)
            item[1].setPosition(pos[0], pos[1], pos[2]);
            item[1].setRotation(pos[3]);
            item[1].setTurretRotation(pos[4]);
            
            if (item[0] == this.self) {
                this.camera.position.x = pos[0];
                this.camera.position.y = pos[1] - 30;
            }
        }

        let items = this.wasm.updateProjectiles(dt);
        for (let i = 0; i < items.length; i += 4) {
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