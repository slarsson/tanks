import * as THREE from 'three';
import { GameMap } from './Map'
import Tank from './Tank';
import Projectile from './Particle';
import Graphics from './Graphics';
import Camera from './Camera';

class Game {

    public camera: Camera;
    private scene: THREE.Scene;
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
       
        // main setup
        this.scene = new THREE.Scene();
        this.camera = new Camera(window.innerWidth, window.innerHeight);
        this.renderer = new THREE.WebGLRenderer({ antialias: true, alpha: false });
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        this.renderer.setClearColor(0x000000, 0);

        let root = document.getElementById('root');
        if (root != null) {
            root.innerHTML = '';
            root.appendChild(this.renderer.domElement);
        } else {
            console.error('GAME: #root element not found :(');
        }

        window.addEventListener('resize', this.resize);

        // object setup
        this.gameMap = new GameMap(this.scene);
        this.players = new Map();
        this.projectiles = new Map();

        // init game
        this.timestamp = performance.now();
        this.animate();
    }

    private resize(evt: Event) {
        const w = evt.target as Window;
        this.camera.setAspect(w.innerWidth,w.innerHeight);
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
                this.camera.setPosition(state[0], state[1] - 40, 30);
                if (this.outOfMap != this.gameMap.outOfMap(state[0], state[1])) {
                    this.outOfMap = !this.outOfMap;
                    if (this.outOfMap) {
                        this.graphics.showOutOfMap(3);
                    } else {
                        this.graphics.hideOutOfMap();
                    }
                }
            }
        }

        let projectiles = this.wasm.updateProjectiles(dt);
        for (let i = 0; i < projectiles.length; i += 5) {
            if (projectiles[i+4] == 0) {
                this.projectiles.get(projectiles[i])?.dispose();
                this.projectiles.delete(projectiles[i]);
                continue
            } 
            
            if (!this.projectiles.has(projectiles[i])) {
                this.projectiles.set(projectiles[i], new Projectile(this.scene));
            }
            this.projectiles.get(projectiles[i])?.set(projectiles[i+1], projectiles[i+2], projectiles[i+3]);
        }

        this.gameMap.update(dt);
        this.camera.update(dt);

        // render
        this.renderer.render(this.scene, this.camera.instance);
        this.timestamp = now;
    }

    addPlayer(id: number): void {
        if (!this.players.has(id)) {
            this.players.set(id, new Tank(this.scene));
        } else {
            console.warn(`GAME: player already exists (id: ${id})`);
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
            console.log('GAME: self player not created yet, try again in 100ms..'); // retard solution ..
            setTimeout(() => this.setSelf(id), 100);
        }
    }

    getSelf(): number {
        return this.self;
    }

    update(data: Float32Array): void {
        for (let i = 0; i < data.length; i += 12) {
            const player = this.players.get(data[i]);
            
            if (player == undefined) {
                this.addPlayer(data[i]);
                return;
            }

            if (data[i+11] == 0 && player.isAlive) {
                player.kill();
            } else if (data[i+11] == 1 && !player.isAlive) {
                player.respawn();
            } 
        }
    }
}

export default Game;