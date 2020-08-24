import * as THREE from 'three';
import { Assets } from './Assets';
import { SingleContainer } from './Container';

class Lifter {

    public mesh: THREE.Group;
    private lines: THREE.Group;
    private liftArm: THREE.Group;
    
    private height: number = 0;
    private maxHeight: number = 12;
    private minHeight: number = 5;

    constructor() {
        this.mesh = new THREE.Group();
        this.lines = new THREE.Group();
        this.liftArm = new THREE.Group();

        {
            let m = new THREE.LineBasicMaterial({color: 0x000000});
            let g1 = new THREE.Geometry();
            g1.vertices.push(
                new THREE.Vector3(1.4, 3.5, this.height), 
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
            this.liftArm.position.z = this.height;
            this.mesh.add(this.liftArm);
        }
        
        let c = new SingleContainer(this.liftArm);
        c.setRotation(0, 0, Math.PI/2);
        c.setPosition(0, 0, -1.875); // RÄKNA RÄTT!!??!!
    }

    private moveToHeight(): void {
        this.liftArm.position.z = this.height;
        for (const child of this.lines.children) {
            // fuck TypeScript, how should this be done?
            if (child instanceof THREE.Line) {
                child.geometry.vertices[0].z = this.height;
                child.geometry.verticesNeedUpdate = true;
            }
        }
    }

    public up(dt: number): boolean {
        this.height += 0.0005 * dt;
        if (this.height > this.maxHeight) {
            this.height = this.maxHeight;
            return false;
        }
        this.moveToHeight();
        return true;
    }

    public down(dt: number): boolean {
        this.height -= 0.0005 * dt;
        if (this.height < this.minHeight) {
            this.height = this.minHeight;
            return false;
        }
        this.moveToHeight();
        return true;
    }

    public setHeight(h: number): void {
        this.height = h;
    }
}


class Crane {

    private scene: THREE.Scene;
    private crane: THREE.Group;
    private test: THREE.Group;
    private lifter: Lifter;

    private dir: boolean = true;
    private sidewaysLength: number = 7;

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
        this.scene.add(this.test);

        // let lifter = new THREE.Group();
        // let lifterMaterial = new THREE.MeshPhongMaterial({color: 0xffff00});
        
        // let lifta = new THREE.Mesh(new THREE.BoxGeometry(3.75, 1, 0.5), lifterMaterial);
        
        // //let liftb = new THREE.LineSegments(new THREE.BoxGeometry(1.5, 6, 0.5), lifterMaterial);

        // let liftb = new THREE.Mesh(new THREE.BoxGeometry(1.5, 6, 0.5), lifterMaterial);
        // let liftc = new THREE.Mesh(new THREE.BoxGeometry(3.75, 1, 0.5), lifterMaterial);
        
        // lifta.position.y = 3.5;
        // liftc.position.y = -3.5;

        // lifter.add(lifta);
        // lifter.add(liftb);
        // lifter.add(liftc);

        // lifter.position.z = 7;

        // this.test.add(lifter);


        // let g1 = new THREE.Geometry();
        // g1.vertices.push(
        //     new THREE.Vector3(1.5, 3.5, 7),
        //     new THREE.Vector3(1.5, 3.5, 14.25)
        // );
        // let line1 = new THREE.Line(g1, new THREE.LineBasicMaterial({color: 0x000000}));

        // let g2 = new THREE.Geometry();
        // g2.vertices.push(
        //     new THREE.Vector3(-1.5, 3.5, 7),
        //     new THREE.Vector3(-1.5, 3.5, 14.25)
        // );
        // let line2 = new THREE.Line(g2, new THREE.LineBasicMaterial({color: 0x000000}));

        // let g3 = new THREE.Geometry();
        // g3.vertices.push(
        //     new THREE.Vector3(1.5, -3.5, 7),
        //     new THREE.Vector3(1.5, -3.5, 14.25)
        // );
        // let line3 = new THREE.Line(g3, new THREE.LineBasicMaterial({color: 0x000000}));

        // let g4 = new THREE.Geometry();
        // g4.vertices.push(
        //     new THREE.Vector3(-1.5, -3.5, 7),
        //     new THREE.Vector3(-1.5, -3.5, 14.25)
        // );
        // let line4 = new THREE.Line(g4, new THREE.LineBasicMaterial({color: 0x000000}));


        // this.test.add(line1);
        // this.test.add(line2);
        // this.test.add(line3);
        // this.test.add(line4);

        // setTimeout(() => {
        //     g1.vertices[0].z = 12;
        //     line1.geometry.verticesNeedUpdate = true;
        // });

        // let swag2 = new THREE.Mesh(
        //     new THREE.BoxGeometry(2, 7, 0.5),
        //     new THREE.MeshPhongMaterial({color: 0xffff00})
        // );

        // swag2.position.z = 7;





        // this.test.add(swag);
        // // this.test.add(swag2);
            
        // this.crane.add(this.test);

        this.scene.add(this.crane);

        this.crane.position.x = -25;
        this.test.position.x = -25;

        // let c = new SingleContainer(this.test);
        // c.setPosition(0, 0, 5);
        // c.setRotation(0, 0, (Math.PI / 2));

        // const crane = Assets.objects.crane.scene;
        //     //crane.rotation.y = Math.random() * 3;
        //     crane.rotation.x = Math.PI / 2;

        //     crane.scale.set(0.5, 0.5, 0.5);
        //     crane.position.x = 20;
            
        //     let child = crane.children[0];
        //     if (child instanceof THREE.Mesh) {
        //         child.material = new THREE.MeshPhongMaterial({color: 0xff0000, side: THREE.DoubleSide}); 
        //     }
    }

    public left(dt: number): void {
        let newp = this.test.position.x - 0.001 * dt;
        if (newp > -this.sidewaysLength) {
            this.test.position.x = newp;
        } else {
            this.test.position.x = -this.sidewaysLength;
        }
    }

    public right(dt: number): void {
        let newp = this.test.position.x + 0.001 * dt;
        if (newp < this.sidewaysLength) {
            this.test.position.x = newp;
        } else {
            this.test.position.x = this.sidewaysLength;
        }
    }

    public update(dt: number): void {
        
        //this.left(dt);

        if (this.dir) {
            //this.left(dt);
            if (!this.lifter.up(dt)) {
                this.dir = false;
            }
        } else {
            //this.right(dt);
            if (!this.lifter.down(dt)) {
                this.dir = true;
            }
        }
    }

}

export default Crane;