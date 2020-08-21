import * as THREE from 'three';
import { Assets } from './Assets';

class Item {
    
    private scene: THREE.Group;
    private mesh: THREE.Mesh;
    
    constructor(scene: THREE.Group) {
        this.scene = scene;
        
        let g = new THREE.BoxGeometry(8, 3.75, 3.75);
        let m = new THREE.MeshPhongMaterial({color: 0x00a0d6});

        let rand = Math.trunc(Math.random()*3);
        let t: THREE.Texture | undefined = undefined;
        if (rand == 0) {
            t = Assets.textures?.postnord.clone();
        } else if (rand == 1) {
            t = Assets.textures?.maersk.clone();
        } else {
            t = Assets.textures?.msc.clone();
        }

        if (t != undefined) {
            t.needsUpdate = true;
            m = new THREE.MeshPhongMaterial({map: t, side: THREE.DoubleSide});
        } else {
            console.warn('CONTAINER: texture is missing :(');
        }
        
        // S2
        g.faceVertexUvs[0][0][0].set(0.32, 0.5);
        g.faceVertexUvs[0][0][1].set(0, 0.5);
        g.faceVertexUvs[0][0][2].set(0.32, 0);

        g.faceVertexUvs[0][1][1].set(0, 0);
        g.faceVertexUvs[0][1][0].set(0, 0.5);
        g.faceVertexUvs[0][1][2].set(0.32, 0);

        // S1
        g.faceVertexUvs[0][2][2].set(0, 1);
        g.faceVertexUvs[0][2][0].set(0, 0.5);
        g.faceVertexUvs[0][2][1].set(0.32, 0.5);
        
        g.faceVertexUvs[0][3][2].set(0, 1);
        g.faceVertexUvs[0][3][1].set(0.32, 1);
        g.faceVertexUvs[0][3][0].set(0.32, 0.5);

        // back
        g.faceVertexUvs[0][4][0].set(1, 0.5);
        g.faceVertexUvs[0][4][1].set(1, 1);
        g.faceVertexUvs[0][4][2].set(0.32, 0.5);

        g.faceVertexUvs[0][5][0].set(1, 1);
        g.faceVertexUvs[0][5][1].set(0.32, 1);
        g.faceVertexUvs[0][5][2].set(0.32, 0.5);

        // front
        g.faceVertexUvs[0][6][0].set(0.32, 1);
        g.faceVertexUvs[0][6][1].set(0.32, 0.5);
        g.faceVertexUvs[0][6][2].set(1, 1);

        g.faceVertexUvs[0][7][0].set(0.32, 0.5);
        g.faceVertexUvs[0][7][1].set(1, 0.5);
        g.faceVertexUvs[0][7][2].set(1, 1);
        
        // top
        g.faceVertexUvs[0][8][0].set(0.32, 0.5);
        g.faceVertexUvs[0][8][2].set(1, 0.5);
        g.faceVertexUvs[0][8][1].set(0.32, 0);

        g.faceVertexUvs[0][9][0].set(0.32, 0);
        g.faceVertexUvs[0][9][2].set(1, 0.5);
        g.faceVertexUvs[0][9][1].set(1, 0);

        // bottom
        g.faceVertexUvs[0][10][0].set(0.32, 0.5);
        g.faceVertexUvs[0][10][2].set(1, 0.5);
        g.faceVertexUvs[0][10][1].set(0.32, 0);

        g.faceVertexUvs[0][11][0].set(0.32, 0);
        g.faceVertexUvs[0][11][2].set(1, 0.5);
        g.faceVertexUvs[0][11][1].set(1, 0);

        this.mesh = new THREE.Mesh(g, m);
        scene.add(this.mesh);
    }

    public setPosition(x: number, y: number, z: number): void {
        this.mesh.position.set(x, y, z);
    }

    public dispose(): void {
        this.scene.remove(this.mesh);
    }
}

class Container {    
    
    static readonly HEIGHT = 3.75;
    static readonly WIDTH = 3.75;
    static readonly LENGTH = 8;

    private group: THREE.Group;
    private containers: Item[] = []; 

    constructor(scene: THREE.Scene, nTotal: number, nBottom: number) {
        if (nBottom > nTotal) {nBottom = nTotal;}

        this.group = new THREE.Group();
        scene.add(this.group);
        
        let yInitValue = -0.5 * Container.WIDTH * nBottom + (Container.WIDTH * 0.5);

        let x = 0;
        let y = yInitValue;
        let z = Container.HEIGHT / 2;
        
        for (let i = 0; i < nTotal; i++) {
            let item = new Item(this.group);
            item.setPosition(x, y, z);
            y += Container.WIDTH;

            if (((i+1) % nBottom) == 0) {
                z += Container.HEIGHT;
                y = yInitValue;
            }

           this.containers.push(item);
        }
    }

    public setPosition(x: number, y: number, z: number): void {
        this.group.position.set(x, y, z);
    }

    public setRotation(rot: number): void {
        this.group.rotation.z = rot;
    }

    public setScale(f: number): void {
        this.group.scale.set(f, f, f);
    }
}

export default Container;