import * as THREE from 'three';

class Tank {

    private scene: THREE.Scene;
    private tank: THREE.Group;
    private hull: THREE.Group;
    private turret: THREE.Group;
    private isAlive: boolean = true;

    constructor(scene: THREE.Scene) {
        this.scene = scene;
        
        this.tank = new THREE.Group();
        this.hull = new THREE.Group();
        this.turret = new THREE.Group();

        let material = new THREE.MeshPhongMaterial({color: 0xc2b280});

        let hullMesh = new THREE.Mesh(new THREE.BoxGeometry(3, 6, 1), material);
            hullMesh.position.z = 0.5;
        
        let turretMesh = new THREE.Mesh(new THREE.BoxGeometry(2, 2, 1), material);
            turretMesh.position.z = 1.5;
 
        let gunMesh = new THREE.Mesh(new THREE.CylinderGeometry(0.2, 0.2, 4, 16), material);
            gunMesh.position.z = 1.5;
            gunMesh.position.y = 2;
 
        this.turret.add(turretMesh);
        this.turret.add(gunMesh);
        this.hull.add(hullMesh);
        this.hull.add(this.turret);
        this.tank.add(this.hull);
        this.scene.add(this.tank);

        // test
        //this.addArrow();
    }

    setPosition(x: number, y: number, z: number) {
        this.tank.position.set(x, y, z);
    } 

    setRotation(rot: number) {
        this.hull.rotation.z = rot;
    }

    setTurretRotation(rot: number) {
        this.turret.rotation.z = rot;
    }

    dispose(): void {
        this.scene.remove(this.tank);
    }

    kill(): void {
        if (this.isAlive) {
            this.isAlive = false;
            this.dispose();
        }
    }

    respawn(): void {
        if (!this.isAlive) {
            this.isAlive = true;
            this.scene.add(this.tank);
        }
    }

    addArrow(rot: number = 0, color: number = 0x00ff00): void {
        let geometry = new THREE.Geometry();
            geometry.vertices= [new THREE.Vector3(0,0,0), new THREE.Vector3(0.5,0,1), new THREE.Vector3(-0.5,0,1)]; 
            geometry.faces = [new THREE.Face3(0, 1, 2)];

        let material = new THREE.MeshBasicMaterial({
            color: color, 
            transparent: true, 
            opacity: 0.5, 
            side: THREE.DoubleSide
        });

        let mesh = new THREE.Mesh(geometry, material);
            mesh.rotation.z = rot;
            mesh.position.z = 3;
            mesh.scale.set(2, 2, 2);

        this.tank.add(mesh);
    }
}

export default Tank;