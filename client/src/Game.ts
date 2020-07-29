import * as THREE from 'three';

import GameMap from './Map'
import Tank from './Tank';
import Projectile from './Particle';
import Graphics from './Graphics';

class Game {

    private scene: THREE.Scene;
    private camera: THREE.PerspectiveCamera;
    private renderer: THREE.WebGLRenderer;

    private timestamp: number;
    private self: number = -1;
    private outOfMap: boolean = false;

    private gameMap: GameMap;
    private players: Map<number, Tank>;
    private projectiles: Map<number, Projectile>;


    private graphics: Graphics;
    private wasm: any;

    constructor(wasm, graphics) {
        this.wasm = wasm;
        this.graphics = graphics;

        this.resize = this.resize.bind(this);
        this.animate = this.animate.bind(this);
       
        this.scene = new THREE.Scene();

        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 1, 2000);
        this.camera.position.y = -50;
        this.camera.position.z = 30;
        this.camera.lookAt(0, 0, 0);

        this.renderer = new THREE.WebGLRenderer({ antialias: true, alpha: false });
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        this.renderer.setClearColor(0x000000, 0);

        document.getElementById('root')?.appendChild(this.renderer.domElement);
        window.addEventListener('resize', this.resize);

        // object setup
        this.gameMap = new GameMap(this.scene);
        this.players = new Map();
        this.projectiles = new Map();

        // init
        this.timestamp = performance.now();
        this.animate();

        // // test
        // this.addPlayer(1337);
        // this.addPlayer(1337);
        // this.addPlayer(1338);
        // this.addPlayer(1339);
        // setTimeout(() => this.removePlayer(1339), 4000);
    }

    private resize(evt: Event) {
        const w = evt.target as Window;
        this.camera.aspect = w.innerWidth / w.innerHeight;
        this.camera.updateProjectionMatrix();
        this.renderer.setSize(w.innerWidth, w.innerHeight);
    }

    private animate(): void {
        requestAnimationFrame(this.animate);

        // ..
        let now = performance.now();
        let dt = now - this.timestamp;

        // update
        for (let { 0: id, 1: player } of this.players[Symbol.iterator]()) {
            const state = this.wasm.getPos(id, dt)
            player.setPosition(state[0], state[1], state[2]);
            player.setRotation(state[3]);
            player.setTurretRotation(state[4]);

            if (id == this.self) {
                this.camera.position.x = state[0];
                this.camera.position.y = state[1] - 40;

                if (this.outOfMap != this.gameMap.outOfMap(state[0], state[1])) {
                    this.outOfMap = !this.outOfMap;
                    
                    if (this.outOfMap) {
                        this.graphics.addMessage('Out of map...', 1000, Graphics.WARNING);
                    } else {
                        this.graphics.addMessage('Welcome back...', 1000, Graphics.INFO);
                    }
                }
            }
        }

        // TODO: fix projectilez
        let projectiles = this.wasm.updateProjectiles(dt);
        for (let i = 0; i < projectiles.length; i += 5) {
            if (projectiles[i+4] == 0) {
                //console.log('remove:', items[i]);
                this.projectiles.get(projectiles[i])?.dispose();
                this.projectiles.delete(projectiles[i]);
                continue
            } 
            
            if (!this.projectiles.has(projectiles[i])) {
                this.projectiles.set(projectiles[i], new Projectile(this.scene));
            }
            this.projectiles.get(projectiles[i])?.set(projectiles[i+1], projectiles[i+2], projectiles[i+3]);
        }

        // render
        this.renderer.render(this.scene, this.camera);
        this.timestamp = now;
    }

    addPlayer(id: number): void {
        if (!this.players.has(id)) {
            this.players.set(id, new Tank(this.scene));
        } else {
            console.warn(`GAME: player alreay exists (id: ${id})`);
        }
    }

    removePlayer(id: number): void {
        this.players.get(id)?.dispose();
        this.players.delete(id);
    }

    setSelf(id: number): void {
        this.self = id;
        if (this.players.has(id)) {
            this.players.get(id)?.addArrow();
        } else {
            console.log('GAME: self player not created yet, try again in 100ms..')
            setTimeout(() => this.setSelf(id), 100);
        }
    }

    getSelf(): number {
        return this.self;
    }

    update(data: Float32Array): void {
        for (let i = 0; i < data.length; i += 12) {
            if (!this.players.has(data[i])) {
                this.addPlayer(data[i]);
            } else if (data[i+9] == 1) {
                console.log('GAME: add projectile for id: ', data[i]);
            } else if (data[i+11] == 0) {
                this.players.get(data[i])?.kill(); // not efficient? :((
            } else {
                this.players.get(data[i])?.respawn(); // meh..
            }  
        }
    }
}

export default Game;