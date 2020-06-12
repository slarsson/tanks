import * as THREE from 'three';
import Vehicle from './Vehicle';
import Keypress from './Keypress';
import Connection from './Connection';

class Render {

    conn: Connection;
    scene: THREE.Scene;
    camera: THREE.PerspectiveCamera;
    renderer: THREE.WebGLRenderer;
    timestamp: number;
    key: Keypress;
    test: Vehicle[];
    state: any;

    constructor() {
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth/window.innerHeight, 1, 2000);
        this.camera.rotation.x = Math.PI / 4;
        //this.camera.rotation.z = Math.PI / 6;
        this.camera.position.y = -10;
        
        this.renderer = new THREE.WebGLRenderer({antialias: true});
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        
        document.getElementById('root')?.appendChild(this.renderer.domElement);

        this.key = new Keypress();

        this.animate = this.animate.bind(this);
        

        // init
        this.timestamp = performance.now();
        
        //this.scene.add(this.cube);
        this.camera.position.z = 15;
        this.test = [
            new Vehicle(this.scene, this.key),
            new Vehicle(this.scene)
        ];
        this.test[1].setColor(0x00ff00);

        this.animate();

        // test websocket
        this.state = [];

        this.conn = new Connection((x: any) => {
            if (x.length == 0) return;
            let last = x[x.length - 1];
            this.test[this.test.length - 1].setPosition(new THREE.Vector3(last.x, last.y, last.z));
            this.state = x;
        });
    }

    animate(): void {
        requestAnimationFrame(this.animate);
        let now = performance.now();
        let dt = now - this.timestamp;

        for (const item of this.test) {
            item.update(dt);
        }
        //this.test.update(dt);
        //this.test.setPosition(new THREE.Vector3(0, 0, 0));

        this.renderer.render(this.scene, this.camera);    
        this.timestamp = now;
    }

}

export default Render;