import * as THREE from 'three';
import { Assets } from './Assets';
import { SingleContainer } from './Container';
import { runInThisContext } from 'vm';

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


interface animation {
    type: number;
    delay: number;
    direction: number;
    speed: number;
    stop: number;
}

class Crane {

    private scene: THREE.Scene;
    private crane: THREE.Group;
    private liftArm: THREE.Group;
    private lifter: Lifter;
    private lifterHeight: number = 12;

    private animation: animation[];
    private animationIndex: number = Infinity;
    private animationDelay: number = 0;

    constructor(scene: THREE.Scene) {
        this.scene = scene;
        this.crane = new THREE.Group();
        this.lifter = new Lifter();
        this.lifter.setHeight(this.lifterHeight);

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
            console.error('CRANE: crane-mesh not found :(');
        }

        this.liftArm = new THREE.Group();

        let block = new THREE.Mesh(
            new THREE.BoxGeometry(3, 11, 1),
            new THREE.MeshPhongMaterial({color: 0xff0000})
        );
        block.position.z = 14.25;
        
        this.liftArm.add(block);
        this.liftArm.add(this.lifter.mesh);
        this.crane.add(this.liftArm);
        this.scene.add(this.crane);

        this.animation = [
            {
                type: 0,
                delay: 0,
                direction: 1,
                speed: 0.01,
                stop: 7.5
            },
            {
                type: 1,
                delay: 500,
                direction: -1,
                speed: 0.01,
                stop: 7.5
            },
            {
                type: 1,
                delay: 200,
                direction: 1,
                speed: 0.01,
                stop: 12
            },
            {
                type: 0,
                delay: 500,
                direction: -1,
                speed: 0.01,
                stop: 0
            },

            // {
            //     type: 0,
            //     delay: 500,
            //     direction: -1,
            //     speed: 0.01,
            //     stop: -7
            // }
        ];

        // this.animation = [
        //     {
        //         type: 'test',
        //         speed: 0.01,
        //         direction: 1,
        //         end: 7,
        //     },
        //     {
        //         type: 'test',
        //         speed: 0.01,
        //         direction: -1,
        //         end: -7,
        //     },
        //     {
        //         type: 'test',
        //         speed: 0.005,
        //         direction: 1,
        //         end: 0,
        //     }
        // ];
    }

    public setPosition(x: number, y: number, z: number): void {
        this.crane.position.set(x, y, z);
    }

    public triggerAnimation(): void {
        console.log('trigge animations..');

        this.animationIndex = 0;
        // this.animationIndex = 0;
        // this.lifterHeight = 12;
        // this.lifter.setHeight(this.lifterHeight);
        // this.liftArm.position.x = 0;
    }

    public update(dt: number): void {
        if (this.animationIndex < this.animation.length) {
            let a = this.animation[this.animationIndex];

            if (a.delay > this.animationDelay) {
                this.animationDelay += dt;
                return;
            }

            if (a.type == 0) {
                this.liftArm.position.x += dt * a.direction * a.speed;
                if ((a.direction == 1 && this.liftArm.position.x > a.stop) || (a.direction == -1 && this.liftArm.position.x < a.stop)) {
                    this.liftArm.position.x = a.stop;
                    this.animationDelay = 0;
                    this.animationIndex++;
                }
            } else if (a.type == 1) {
                this.lifterHeight += dt * a.direction * a.speed;
                if ((a.direction == 1 && this.lifterHeight > a.stop) || (a.direction == -1 && this.lifterHeight < a.stop)) {
                    this.lifterHeight = a.stop;
                    this.animationDelay = 0;
                    this.animationIndex++;
                }
                this.lifter.setHeight(this.lifterHeight);
            }

            //console.log(a);
            
            // if ((a.direction == 1 && this.liftArm.position.x > a.stop) || (a.direction == -1 && this.liftArm.position.x < a.stop)) {
            //     if (a.type == 0) {
            //         this.liftArm.position.x = a.stop;
            //     } else if (a.type == 1) {
            //         this.lifterHeight = a.stop;
            //         this.lifter.setHeight(this.lifterHeight);
            //     }
            //     this.animationDelay = 0;
            //     this.animationIndex++;
            // }

            // if ((a.direction == 1 && this.liftArm.position.x > a.stop) || (a.direction == -1 && this.liftArm.position.x < a.stop)) {
            //     if (a.type == 0) {
            //         this.liftArm.position.x = a.stop;
            //     } else if (a.type == 1) {
            //         this.lifterHeight = a.stop;
            //         this.lifter.setHeight(this.lifterHeight);
            //     }
            //     this.animationDelay = 0;
            //     this.animationIndex++;
            // }

        }


        // let a = this.animation[this.animationIndex];
        // if (a != undefined) {
        //     this.liftArm.position.x += dt * a.direction * a.speed;
        //     if (
        //         (a.direction == 1 && this.liftArm.position.x > a.end) || 
        //         (a.direction == -1 && this.liftArm.position.x < a.end)
        //     ) {
        //         this.liftArm.position.x = a.end;
        //         this.animationIndex++;
        //     }
        // }
    }
}

export default Crane;