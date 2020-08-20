import * as THREE from 'three';

class Camera {
    
    public static readonly LOBBY: number = 0;
    public static readonly FOLLOW: number = 1;
    public static readonly FLY: number = 2;

    public instance: THREE.PerspectiveCamera;
    private mode: number = -1;
    private matrix:  THREE.Matrix4;
    private flyState = {
        speed: {x: 0, y: 0, z: 0},
        stop: {x: 0, y: 0, z: 0},
        target: {x: 0, y: 0, z: 0}
    };

    constructor(w: number, h: number) {
        this.instance = new THREE.PerspectiveCamera(75, (w / h), 1, 2000);
        this.matrix = new THREE.Matrix4();
        this.setMode(Camera.LOBBY);
    }

    setPosition(x: number, y: number, z: number): void {
        if (this.mode != Camera.FOLLOW) return;
        this.instance.position.set(x, y, z);
    }

    setAspect(w: number, h: number): void {
        this.instance.aspect = w / h;
        this.instance.updateProjectionMatrix();
    }

    setFlyTo(x: number, y: number, z: number, t: number): void {
        this.flyState.stop.x = x;
        this.flyState.speed.x = (this.flyState.stop.x - this.instance.position.x) / t;
        this.flyState.stop.y = y;
        this.flyState.speed.y = (this.flyState.stop.y - this.instance.position.y) / t;
        this.flyState.stop.z = z;
        this.flyState.speed.z = (this.flyState.stop.z - this.instance.position.z) / t;
    }

    setFlyToTarget(x: number, y: number, z: number): void {
        this.flyState.target.x = x;
        this.flyState.target.y = y;
        this.flyState.target.z = z;
    }

    setMode(mode: number, x?: number, y?:number, z?: number): void {
        if (mode != Camera.LOBBY && mode != Camera.FOLLOW && mode != Camera.FLY) return;
        
        if (mode == Camera.FOLLOW) {
            this.instance.position.set(x != undefined ? x : 0, y != undefined ? y : -40, z != undefined ? z : 30);
            this.instance.lookAt(0, 0, 0);
        } else if (mode == Camera.LOBBY) {
            this.instance.position.set(x != undefined ? x : 0, y != undefined ? y : -60, z != undefined ? z : 30);
            this.instance.lookAt(0, 0, 0);
        }
        
        this.mode = mode;
        console.log(`CAMERA: mode set to => ${this.mode}`);
    }

    update(dt: number): void {
        if (this.mode == Camera.FOLLOW) return;
        
        if (this.mode == Camera.LOBBY) {
            this.matrix.makeRotationZ(dt * 0.0001);
            this.instance.applyMatrix4(this.matrix);
            return;
        }

        if (this.mode == Camera.FLY) {
            this.instance.position.x += this.flyState.speed.x * dt;
            this.instance.position.y += this.flyState.speed.y * dt;
            this.instance.position.z += this.flyState.speed.z * dt;

            this.instance.up.set(0, 0, 1);
            this.instance.lookAt(this.flyState.target.x, this.flyState.target.y, this.flyState.target.z);

            // only have to check one of x, y or z
            if (this.flyState.speed.x >= 0) {
                if (this.instance.position.x > this.flyState.stop.x) {
                    this.setMode(Camera.FOLLOW);
                }
            } else {
                if (this.instance.position.x < this.flyState.stop.x) {
                    this.setMode(Camera.FOLLOW);
                }
            } 
            return;          
        }
    }
}

export default Camera;