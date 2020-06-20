import * as THREE from 'three';
import { Colors } from 'three';
import Keypress from './Keypress';

class Vehicle {

    protected body: THREE.Mesh;
    protected geometry: THREE.Geometry;
    protected material: THREE.MeshBasicMaterial;
    static readonly COLORS = [0xff0000, 0x00ff00, 0x0000ff]; 
    private keys: Keypress;


    constructor(scene: THREE.Scene, keys: Keypress) {
        this.keys = keys;
        
        this.geometry = new THREE.BoxGeometry(2, 4, 1);
        this.material = new THREE.MeshBasicMaterial({color: Vehicle.COLORS[Math.floor(Math.random() * Vehicle.COLORS.length)]});
        this.body = new THREE.Mesh(this.geometry, this.material);

        scene.add(this.body);

        this.setPosition = this.setPosition.bind(this);
        this.setColor = this.setColor.bind(this);
        this.update = this.update.bind(this);
    }

    setPosition(v: THREE.Vector3) {
        this.body.position.set(v.x, v.y, v.z);
    } 

    setColor(color: number) {
        this.material.setValues({color: color});
    }

    update(dt: number): void {
        let x = 0.005 * dt;
        if (this.keys.status.w) this.body.position.y += x;
        if (this.keys.status.s) this.body.position.y -= x;
        if (this.keys.status.a) this.body.position.x -= x;
        if (this.keys.status.d) this.body.position.x += x; 
        //console.log('do update..');
    } 
}

export default Vehicle;