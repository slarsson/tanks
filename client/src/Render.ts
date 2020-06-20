import * as THREE from 'three';
import Vehicle from './Vehicle';
import TestVechicle from './TestVehicle';
import Keypress from './Keypress';
import Connection from './Connection';

class Render {

    conn: Connection;
    scene: THREE.Scene;
    camera: THREE.PerspectiveCamera;
    renderer: THREE.WebGLRenderer;
    timestamp: number;
    key: Keypress;
    
    vehicles: Map<number, Vehicle>;
    //player: TestVechicle;

    wasm: any;

    constructor(wasm: any) {
        this.wasm = wasm;
        this.wasm.test(1, 2, 3);
        
        // setup
        this.scene = new THREE.Scene();

        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth/window.innerHeight, 1, 2000);
        this.camera.rotation.x = Math.PI / 4;
        this.camera.position.y = -10;
        this.camera.position.z = 15;
        
        this.renderer = new THREE.WebGLRenderer({antialias: true});
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        document.getElementById('root')?.appendChild(this.renderer.domElement);

        this.key = new Keypress();

        // method binds
        this.animate = this.animate.bind(this);
        this.test = this.test.bind(this);
        this.worldTick = this.worldTick.bind(this);
        
        // test
        //this.player = new TestVechicle(this.scene, this.key);
        this.vehicles = new Map();

        // websocket
        this.conn = new Connection(this.worldTick
            // (x: any) => this.wasm.update(...x) 
            
            //(state: any) => {
            //console.log(state);
            // //console.log(state);
            // if (state.length == 0) return;
            // if (state.length != this.vehicles.length) {
            //     this.vehicles.push(new Vehicle(this.scene));
            // }

            // for (let i = 0; i < state.length; i++) {
            //     this.vehicles[i].setPosition(new THREE.Vector3(state[i].x, state[i].y, state[i].z));
            // }
            ,
            () => {
                this.test();
            }
        );


        window.addEventListener('keydown', (evt: KeyboardEvent) => {
            //console.log(evt.key);
            this.wasm.state(evt.key, true);
        });

        window.addEventListener('keyup', (evt: KeyboardEvent) => {
            //console.log(evt.key);
            this.wasm.state(evt.key, false);
        });

        // init
        this.timestamp = performance.now();
        this.animate();
        //this.test();

 
        //window.addEventListener('keyup', this.remove);
    }

    worldTick(state: Float32Array): void {
        //console.log(state);
        this.wasm.update(...state);
    
        for (let i = 0; i < state.length; i += 4) {
            if (this.vehicles.has(state[i])) {
                const pos = this.wasm.getPos(state[i])
                //console.log(pos);
                
                this.vehicles.get(state[i])?.setPosition(new THREE.Vector3(pos[0], pos[1], pos[2]));


                //this.vehicles.get(state[i])?.setPosition(new THREE.Vector3(state[i+1], state[i+2], state[i+3]));
            } else {
                this.vehicles.set(state[i], new Vehicle(this.scene, this.key));
            }

            //console.log(state[i], state[i+1], state[i+2]);
        }
    }

    test(): void {
        //this.wasm.print(); 

        //console.log(this.wasm.poll())

        this.conn.send(this.wasm.poll().buffer);
        // this.conn.send(this.key.poll().buffer);
        setTimeout(this.test, 50);
    }

    animate(): void {
        requestAnimationFrame(this.animate);
        let now = performance.now();
        let dt = now - this.timestamp;

       
        // this.wasm.local(dt)

        // const it = this.vehicles[Symbol.iterator]();
        // for (let item of it) {
        //     const pos = this.wasm.getPos(item[0])
        //     item[1].setPosition(new THREE.Vector3(pos[0], pos[1], pos[2]));
        // }


       

        this.renderer.render(this.scene, this.camera);    
        this.timestamp = now;
    }

}

export default Render;