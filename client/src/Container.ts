import * as THREE from 'three';
import { Assets } from './AssetsTest';
import { Vector2 } from 'three';

class Item {
    
    private static readonly COLORS = [
        0xff6600, // Hapag-Lloyd
        0x00a0d6, // POSTNORD
        0x04246A, // CMA CGM
        0xa6de7b, // CHINA SHIPPING
        0xbd4805, // unknown
        0xf1ece1, // MEARSK
    ];

    private scene: THREE.Group;
    private mesh: THREE.Mesh;
    
    constructor(scene: THREE.Group) {
        this.scene = scene;
        
        let g = new THREE.BoxGeometry(8, 3.75, 3.75);
        
        let m = new THREE.MeshPhongMaterial({color: 0x00a0d6});

        
        let swag = Math.trunc(Math.random()*3);
        let t: THREE.Texture | undefined = undefined;
        if (swag == 0) {t = Assets.textures?.postnord.clone();}
        if (swag == 1) {t = Assets.textures?.msc.clone();}
        if (swag == 2) {t = Assets.textures?.maersk.clone();}

        // let t = Assets.textures?.msc.clone();
        if (t != undefined) {
            t.needsUpdate = true;
            m = new THREE.MeshPhongMaterial({map: t, side: THREE.DoubleSide});
        }

        let zero = [new Vector2(0, 0), new Vector2(0, 0), new Vector2(0, 0)]; 
        g.faceVertexUvs[0][0] = zero;
        g.faceVertexUvs[0][1] = zero;
        g.faceVertexUvs[0][2] = zero;
        g.faceVertexUvs[0][3] = zero;
        g.faceVertexUvs[0][8] = zero;
        g.faceVertexUvs[0][9] = zero;
        g.faceVertexUvs[0][10] = zero;
        g.faceVertexUvs[0][11] = zero;

        let st = [
            new Vector2(0.32, 0),
            new Vector2(1, 0),
            new Vector2(1, 1),
            new Vector2(0.32, 1)
        ];

        g.faceVertexUvs[0][4] = [st[1], st[2], st[0]];
        g.faceVertexUvs[0][5] = [st[2], st[3], st[0]];

        g.faceVertexUvs[0][6] = [st[3], st[0], st[2]];
        g.faceVertexUvs[0][7] = [st[0], st[1], st[2]];

        
        // front
        // g.faceVertexUvs[0][6][0].set(0.32, 1);
        // g.faceVertexUvs[0][6][1].set(0.32, 0);
        // g.faceVertexUvs[0][6][2].set(1, 1);

        // g.faceVertexUvs[0][7][0].set(0.32, 0);
        // g.faceVertexUvs[0][7][1].set(1, 0);
        // g.faceVertexUvs[0][7][2].set(1, 1);

        // // back
        // g.faceVertexUvs[0][4][0].set(1, 0);
        // g.faceVertexUvs[0][4][1].set(1, 1);
        // g.faceVertexUvs[0][4][2].set(0.32, 0);

        // g.faceVertexUvs[0][5][0].set(1, 1);
        // g.faceVertexUvs[0][5][1].set(0.32, 1);
        // g.faceVertexUvs[0][5][2].set(0.32, 0);


        this.mesh = new THREE.Mesh(
            g,
            m
        );
            
        //console.log(g.faceVertexUvs);

        scene.add(this.mesh);

        //this.mesh.rotation.z = 0.2;
        // let f = () => {this.mesh.rotation.z += 0.1; setTimeout(() => f(), 50)};
        // f();
    }

    public setPosition(x: number, y: number, z: number): void {
        this.mesh.position.set(x, y, z);
    }

    public dispose(): void {
        this.scene.remove(this.mesh);
    }

    public test(tex: THREE.Texture): void {
        console.log('wtf..', tex);

        this.mesh.material = [
            new THREE.MeshPhongMaterial({color: 0xff0000}),
            new THREE.MeshPhongMaterial({color: 0xff0000}),
            new THREE.MeshPhongMaterial({color: 0xff0000}),
            new THREE.MeshPhongMaterial({map: tex, color: 0xff0000}),
            new THREE.MeshPhongMaterial({color: 0xff0000}),
            new THREE.MeshPhongMaterial({color: 0xff0000})
        ];
    }
}

class Container {    
    
    static readonly HEIGHT = 3.75;
    static readonly WIDTH = 3.75;
    static readonly LENGTH = 8;

    private group: THREE.Group;

    private containers: Item[] = []; 

    constructor(scene: THREE.Scene, nTotal: number, nBottom: number) {
        if (nBottom > nTotal) {
            nBottom = nTotal;
        }
        
        this.group = new THREE.Group();
        //this.group.rotation.z = 0.2;
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
        //this.group.rotation.z = Math.random() * Math.PI * 2;
    }

    public setRotation(rot: number): void {
        this.group.rotation.z = rot;
    }

    public test(tex: THREE.Texture): void {
        this.containers[0].test(tex);
        // let target = this.group.children[0];
        // console.log(target.material);
    }
}

export default Container;