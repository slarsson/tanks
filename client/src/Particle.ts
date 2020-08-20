import * as THREE from 'three';

class Projectile {

    private mesh: THREE.Mesh;
    private scene: THREE.Scene;

    constructor(scene: THREE.Scene) {
        this.mesh = new THREE.Mesh(
            new THREE.SphereGeometry(0.15, 32, 32), 
            new THREE.MeshBasicMaterial({color: 0xffffff})
        );

        this.scene = scene;
        this.scene.add(this.mesh);
    }

    dispose(): void {
        this.scene.remove(this.mesh);
    }

    set(x: number, y: number, z: number): void {
        this.mesh.position.x = x
        this.mesh.position.y = y
        this.mesh.position.z = z
    }
}

export default Projectile;