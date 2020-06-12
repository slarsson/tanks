import * as THREE from 'three';
import Keypress from './Keypress';

class Vehicle {

    body: THREE.Mesh;
    keys: Keypress | null;

    private material: THREE.MeshBasicMaterial;

    constructor(scene: THREE.Scene, keys?: Keypress) {
        
        if (keys) {
            this.keys = keys;
        } else {
            this.keys = null;
        }

        this.material = new THREE.MeshBasicMaterial({color: 0xff0000});
        this.body = new THREE.Mesh(
            new THREE.BoxGeometry(2, 4, 1),
            this.material
        );

        scene.add(this.body);    

        this.update = this.update.bind(this);
        this.setPosition = this.setPosition.bind(this);
    }

    update(dt: number): void {
        if (this.keys?.status.w) this.body.position.y += 0.005 * dt;
        if (this.keys?.status.s) this.body.position.y -= 0.005 * dt;
        if (this.keys?.status.a) this.body.position.x -= 0.005 * dt;
        if (this.keys?.status.d) this.body.position.x += 0.005 * dt; 
    }

    setPosition(v: THREE.Vector3) {
        this.body.position.set(v.x, v.y, v.z);
        //console.log('do shiet');
    } 

    setColor(color: number) {
        this.material.setValues({color: color});
    }
}

export default Vehicle;