import * as THREE from 'three';
import { Colors } from 'three';
import Keypress from './Keypress';

class Vehicle {

    protected body: THREE.Group;
    protected geometry: THREE.Geometry;
    protected material: THREE.MeshBasicMaterial;
    static readonly COLORS = [0xff0000, 0x00ff00, 0x0000ff]; 
    //private keys: Keypress;

    private scene: THREE.Scene;


    protected turret: THREE.Group;
    protected turretMesh: THREE.Mesh;
    protected hull: THREE.Mesh;
    protected gun: THREE.Mesh;


    constructor(_scene: THREE.Scene) {
        //this.keys = keys;
        
        this.geometry = new THREE.BoxGeometry(2, 4, 1);
        this.material = new THREE.MeshPhongMaterial({color: 0xc2b280});
        
        
        this.body = new THREE.Group();
        this.turret = new THREE.Group();
        //this.body = new THREE.Mesh(this.geometry, this.material);

        this.hull = new THREE.Mesh(
            new THREE.BoxGeometry(3, 6, 1),
            this.material
        );

        this.turretMesh = new THREE.Mesh(
            new THREE.BoxGeometry(2, 2, 1),
            new THREE.MeshPhongMaterial({color: 0xc2b280})
        );
        this.turretMesh.position.z = 1;

        this.gun = new THREE.Mesh(
            new THREE.CylinderGeometry(0.2, 0.2, 4, 16),
            new THREE.MeshPhongMaterial({color: 0xc2b280})
        );
        this.gun.position.z = 1;
        this.gun.position.y = 2;

        this.body.add(this.hull);
        this.body.add(this.turret);
        
        this.turret.add(this.gun);
        this.turret.add(this.turretMesh);



        this.scene = _scene;
        this.scene.add(this.body);

        this.setPosition = this.setPosition.bind(this);
        this.setColor = this.setColor.bind(this);
        this.update = this.update.bind(this);
    }

    setPosition(x: number, y: number, z: number) {
        this.body.position.set(x, y, z);
    } 

    setRotation(rot: number) {
        this.body.rotation.z = rot;
    }

    setTurretRotation(rot: number) {
        this.turret.rotation.z = rot;
    }

    setColor(color: number) {
        this.material.setValues({color: color});
    }

    getGunOrigin(): THREE.Vector3 {
        console.log(this.gun.position.x);
        return this.body.position.clone();
    }

    getGunRotation(): number {
        return this.turret.rotation.z;
    }

    update(dt: number): void {
        let x = 0.005 * dt;
        // if (this.keys.status.w) this.body.position.y += x;
        // if (this.keys.status.s) this.body.position.y -= x;
        // if (this.keys.status.a) this.body.position.x -= x;
        // if (this.keys.status.d) this.body.position.x += x; 
        //console.log('do update..');
    } 

    dispose(): void {
        this.scene.remove(this.body);
    }
}

export default Vehicle;