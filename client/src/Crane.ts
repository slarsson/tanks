import * as THREE from 'three';
import { Assets } from './Assets';
import { SingleContainer } from './Container';

class Lifter {

    public mesh: THREE.Group;
    public liftArm: THREE.Group;
    private lines: THREE.Group;
    private maxHeight: number = 12;

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

enum AnimationType {
    Horizontal,
    Vertical,
    Trigger
}

interface animation {
    type: AnimationType;
    delay: number;
    direction: number; // TRIGGER = add or remove container, otherwise up/down or left/right
    speed: number;
    stop: number; // TRIGGER = taken index, otherwise end position
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

    private lifterPositions: number[] = [7.5, 3.75, 0, -3.75, -7.5];
    private containerPositions: number[] = [-17.5, -21.25, -25, -28.75, -32.5];
    private containerGroupTop: Array<SingleContainer | null>;
    private containerGroupBottom: Array<SingleContainer | null>;
    private current: SingleContainer | null;

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

        // animations + containers
        this.animation = [
            {type: AnimationType.Horizontal, delay: 200, direction: 0, speed: 0.01, stop: 0},
            {type: AnimationType.Vertical, delay: 300, direction: -1, speed: 0.01, stop: 7.5},
            {type: AnimationType.Trigger, delay: 0, direction: 0, speed: 0, stop: 0},
            {type: AnimationType.Vertical, delay: 300, direction: 1, speed: 0.01, stop: 12},
            {type: AnimationType.Horizontal, delay: 200, direction: 0, speed: 0.01, stop: 0},
            {type: AnimationType.Vertical, delay: 300, direction: -1, speed: 0.01, stop: 7.5},
            {type: AnimationType.Trigger, delay: 0, direction: 0, speed: 0, stop: 0},
            {type: AnimationType.Vertical, delay: 300, direction: 1, speed: 0.01, stop: 12},
            {type: AnimationType.Horizontal, delay: 200, direction: 0, speed: 0.01, stop: 0}
        ];

        this.current = new SingleContainer(this.lifter.liftArm);
        this.current.setRotation(0, 0, Math.PI/2);
        this.current.setPosition(0, 0, -1.875);

        this.containerGroupTop = [null, null, null, null, null];
        this.containerGroupBottom = [null, null, null, null, null];

        for (let i = 0; i < (Math.trunc(Math.random() * 4) + 1); i++) {
            let c = new SingleContainer(this.scene);
            c.setRotation(0, 0, Math.PI/2);
            c.setPosition(this.containerPositions[i], 0, 5.625);
            this.containerGroupTop[i] = c;
        }

        for (let i = (Math.trunc(Math.random() * 4) + 1); i > 0; i--) {
            let c = new SingleContainer(this.scene);
            c.setRotation(0, 0, Math.PI/2);
            c.setPosition(this.containerPositions[5 - i], -35, 5.625);
            this.containerGroupBottom[5 - i] = c;
        }
    }

    public setPosition(x: number, y: number, z: number): void {
        this.crane.position.set(x, y, z);
    }

    public triggerAnimation(): void {
        let target = this.crane.position.y < -17 ? this.containerGroupBottom : this.containerGroupTop;
        
        let free: number[] = [];
        let taken: number[] = []; 
        for (let i = 0; i < target.length; i++) {
            if (target[i] != null) {
                taken.push(i);
            } else {
                free.push(i);
            }
        }

        let f = free[Math.trunc(Math.random() * free.length)];
        let t = taken[Math.trunc(Math.random() * taken.length)];

        // 0, 4, 8
        this.animation[0].stop = this.lifterPositions[f];
        this.animation[0].direction = this.lifterPositions[f] < 0 ? -1 : 1;
        this.animation[4].stop = this.lifterPositions[t];
        this.animation[4].direction = t < f ? 1 : -1;
        this.animation[8].direction = t < 2 ? -1 : 1;

        // container
        this.animation[2].direction = -1; // release
        this.animation[2].stop = f; // index
        this.animation[6].direction = 1; // release
        this.animation[6].stop = t; // index

        // start anime
        this.animationIndex = 0;
    }

    public update(dt: number): void {
        if (this.animationIndex < this.animation.length) {            
            let a = this.animation[this.animationIndex];
            if (a.delay > this.animationDelay) {
                this.animationDelay += dt;
                return;
            }

            if (a.type == AnimationType.Horizontal) {
                this.liftArm.position.x += dt * a.direction * a.speed;
                if ((a.direction == 1 && this.liftArm.position.x > a.stop) || (a.direction == -1 && this.liftArm.position.x < a.stop)) {
                    this.liftArm.position.x = a.stop;
                    this.animationDelay = 0;
                    this.animationIndex++;
                }
            } else if (a.type == AnimationType.Vertical) {
                this.lifterHeight += dt * a.direction * a.speed;
                if ((a.direction == 1 && this.lifterHeight > a.stop) || (a.direction == -1 && this.lifterHeight < a.stop)) {
                    this.lifterHeight = a.stop;
                    this.animationDelay = 0;
                    this.animationIndex++;
                }
                this.lifter.setHeight(this.lifterHeight);
            } else if (a.type == AnimationType.Trigger) {
                let target = this.containerGroupTop;
                let y = 0;
                if (this.crane.position.y < -17) {
                    target = this.containerGroupBottom;
                    y = -35;
                }

                if (a.direction == 1) {
                    this.current = target[a.stop];
                    target[a.stop] = null;
                    if (this.current != null) {
                        this.current.addTo(this.lifter.liftArm);
                        this.current.setPosition(0, 0, -1.875);
                    }
                } else if (a.direction == -1) {
                    if (this.current != null) {
                        let c = this.current;
                        target[a.stop] = c;
                        this.current = null;

                        c.addTo(this.scene);
                        c.setPosition(this.containerPositions[a.stop], y, 1.5*3.75);
                    }
                }

                this.animationDelay = 0;
                this.animationIndex++;
            }
        }
    }
}

export default Crane;