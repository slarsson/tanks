import * as THREE from 'three';
import Vehicle from './Vehicle';
import Keypress from './Keypress';

class TestVechicle extends Vehicle {
    
    //private keys: Keypress;

    constructor(scene: THREE.Scene, keys: Keypress) {
        super(scene, keys);

        //this.keys = keys;

        this.updatez = this.update.bind(this);
    }

    updatez(dt: number): void {
        // let x = 0.005 * dt;
        // if (this.keys.status.w) this.body.position.y += x;
        // if (this.keys.status.s) this.body.position.y -= x;
        // if (this.keys.status.a) this.body.position.x -= x;
        // if (this.keys.status.d) this.body.position.x += x; 
    }


}

export default TestVechicle;