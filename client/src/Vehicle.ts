import * as THREE from 'three';
import { Colors } from 'three';

class Vehicle {

    protected body: THREE.Mesh;
    protected geometry: THREE.Geometry;
    protected material: THREE.MeshBasicMaterial;
    static readonly COLORS = [0xff0000, 0x00ff00, 0x0000ff]; 
    
    constructor(scene: THREE.Scene) {
        this.geometry = new THREE.BoxGeometry(2, 4, 1);
        this.material = new THREE.MeshBasicMaterial({color: Vehicle.COLORS[Math.floor(Math.random() * Vehicle.COLORS.length)]});
        this.body = new THREE.Mesh(this.geometry, this.material);

        scene.add(this.body);

        this.setPosition = this.setPosition.bind(this);
        this.setColor = this.setColor.bind(this);
    }

    setPosition(v: THREE.Vector3) {
        this.body.position.set(v.x, v.y, v.z);
    } 

    setColor(color: number) {
        this.material.setValues({color: color});
    }
}

export default Vehicle;