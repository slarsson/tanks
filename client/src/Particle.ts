import * as THREE from 'three';

class Projectile {

    private mesh: THREE.Mesh;
    private scene: THREE.Scene;

    constructor(scene: THREE.Scene) {
        this.mesh = new THREE.Mesh(
            new THREE.SphereGeometry(0.2, 32, 32), 
            new THREE.MeshBasicMaterial({color: 0xffffff})
        );

        this.scene = scene;
        this.scene.add(this.mesh);
    }

    set(x: number, y: number, z: number): void {
        //console.log(x, y);
        this.mesh.position.x = x
        this.mesh.position.y = y
        this.mesh.position.z = z
    }
}

// class Projectile {
    
//     private mesh: THREE.Mesh;
//     private scene: THREE.Scene;
//     private f: number;
//     private decay: number;

//     private swag: number;

//     constructor(scene: THREE.Scene, rot: number) {
//         this.mesh = new THREE.Mesh(
//             new THREE.SphereGeometry(0.2, 32, 32), 
//             new THREE.MeshBasicMaterial({color: 0xffffff})
//         );

//         this.swag = rot;
//         //this.mesh.position.copy(v);

//         this.scene = scene;
//         this.scene.add(this.mesh);
//         this.f = 0.05;
//         this.decay = 0.001;
//     }

//     update(dt: number) {
//         this.mesh.position.y += this.f * Math.cos(this.swag) * dt;
//         this.mesh.position.x -= this.f * Math.sin(this.swag) * dt;

//         this.mesh.position.z -= this.decay * dt;

//         if (this.mesh.position.y < -1) {
//             this.dispose();
//         }
//     } 

//     dispose() {
//         this.scene.remove(this.mesh);
//     }
// }

// class Particle {
    
//     private particles: Array<Projectile>;
//     private scene: THREE.Scene;

//     constructor(scene: THREE.Scene) {
//         this.scene = scene;
//         this.particles = [];

//         // for (let i = 0; i < 10; i++) {
//         //     this.particles.push(new Projectile(this.scene));
//         // }
//     }

//     add(v: THREE.Vector3): void {
//         //this.particles.push(new Projectile(this.scene, v));
//     }

//     add2(rot: number): void {
//         this.particles.push(new Projectile(this.scene, rot));
//     }

//     update(dt: number): void {
//         for (const p of this.particles) {
//             p.update(dt);
//         }
//         // this.p.position.y += 0.001 * dt;
//         // this.p.position.x += 0.002 * dt;
//     }
// }

export default Projectile;