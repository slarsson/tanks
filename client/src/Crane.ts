import * as THREE from 'three';
import { Assets } from './Assets';
import { SingleContainer } from './Container';

class Lifter {

    public mesh: THREE.Group;
    private lines: THREE.Group;
    private liftArm: THREE.Group;
    
    private maxHeight: number = 12;
    private minHeight: number = 3.75;

    constructor() {
        this.mesh = new THREE.Group();
        this.lines = new THREE.Group();
        this.liftArm = new THREE.Group();

        {
            let m = new THREE.LineBasicMaterial({color: 0x000000});
            let g1 = new THREE.Geometry();
            g1.vertices.push(
                new THREE.Vector3(1.4, 3.5, this.maxHeight), 
                new THREE.Vector3(1.4, 3.5, this.maxHeight+2)
            );
            this.lines.add(new THREE.Line(g1, m));
    
            let g2 = g1.clone();
            g2.vertices[0].x *= -1;
            g2.vertices[1].x *= -1;
            this.lines.add(new THREE.Line(g2, m));
    
            let g3 = g1.clone();
            g3.vertices[0].y *= -1;
            g3.vertices[1].y *= -1;
            this.lines.add(new THREE.Line(g3, m));
    
            let g4 = g3.clone();
            g4.vertices[0].x *= -1;
            g4.vertices[1].x *= -1;
            this.lines.add(new THREE.Line(g4, m));
        
            this.mesh.add(this.lines);
        }

        {
            let m = new THREE.MeshPhongMaterial({color: 0xffff00});
            let top = new THREE.Mesh(new THREE.BoxGeometry(3.75, 1, 0.5), m);
            let middle = new THREE.Mesh(new THREE.BoxGeometry(1.5, 6, 0.5), m);
            let bottom = new THREE.Mesh(new THREE.BoxGeometry(3.75, 1, 0.5), m);

            top.position.y = -3.5;
            top.position.z = 0.25;
            middle.position.z = 0.25;
            bottom.position.y = 3.5;
            bottom.position.z = 0.25;

            this.liftArm.add(top);
            this.liftArm.add(middle);
            this.liftArm.add(bottom);
            this.liftArm.position.z = this.maxHeight;
            this.mesh.add(this.liftArm);
        }
        
        let c = new SingleContainer(this.liftArm);
        c.setRotation(0, 0, Math.PI/2);
        c.setPosition(0, 0, -1.875); // RÄKNA RÄTT!!??!!
    }

    public setHeight(h: number): void {
        this.liftArm.position.z = h;
        for (const child of this.lines.children) {
            // fuck TypeScript, how should this be done?
            if (child instanceof THREE.Line) {
                child.geometry.vertices[0].z = h;
                child.geometry.verticesNeedUpdate = true;
            }
        }
    }
}

class Crane {

    private scene: THREE.Scene;
    private crane: THREE.Group;
    private test: THREE.Group;
    private lifter: Lifter;

    private animationTime: number = 0;
    private animationRunning: boolean = false;
    private animationTest: number = 1;
    private height: number = 12;

    constructor(scene: THREE.Scene) {
        this.scene = scene;
        this.crane = new THREE.Group();

        let mesh = Assets.objects?.crane?.scene.clone();
        if (mesh != undefined) {
            let child = mesh.children[0];
            if (child instanceof THREE.Mesh) {
                child.material = new THREE.MeshPhongMaterial({color: 0xff0000}); 
            }
            mesh.scale.set(0.5, 0.5, 0.5);
            mesh.rotation.x = Math.PI / 2;
            this.crane.add(mesh);
        } else {
            console.error('CRANE: could not found crane-mesh');
        }

        this.test = new THREE.Group();

        let swag = new THREE.Mesh(
            new THREE.BoxGeometry(3, 11, 1),
            new THREE.MeshPhongMaterial({color: 0xff0000})
        );

        swag.position.z = 14.25;
        this.lifter = new Lifter();
        //this.scene.add(this.lifter.mesh);

        this.test.add(swag);
        this.test.add(this.lifter.mesh);
            
        this.crane.add(this.test);
        

        this.scene.add(this.crane);

    }

    public setPosition(x: number, y: number, z: number): void {
        if (!this.animationRunning && y == this.crane.position.y) {
            this.animationTime = 0;
            this.animationRunning = true;
        }
        
        this.crane.position.set(x, y, z);
        //this.test.position.set(x, y, z);
    }

    public update(dt: number): void {
        if (this.animationRunning) {
            this.animationTime += dt;

            this.height -= this.animationTest * dt * 0.001;
            this.lifter.setHeight(this.height);
            this.test.position.x += this.animationTest * dt * 0.002;

            if (this.animationTime >= 3200) {
                this.animationTime = 0;
                this.animationRunning = false;
                this.animationTest *= -1;
            }
        }
    }

}

export default Crane;