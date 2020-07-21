import * as THREE from 'three';

class Item {
    
    private static readonly COLORS = [
        0xff6600, // Hapag-Lloyd
        0x01a0d7, // POSTNORD
        0x04246A, // CMA CGM
        0xa6de7b, // CHINA SHIPPING
        0xbd4805, // unknown
        0xf1ece1, // MEARSK
    ];

    private scene: THREE.Group;
    private mesh: THREE.Mesh;
    
    constructor(scene: THREE.Group) {
        this.scene = scene;
        this.mesh = new THREE.Mesh(
            new THREE.BoxGeometry(8, 3.75, 3.75),
            new THREE.MeshPhongMaterial({color: Item.COLORS[Math.trunc(Math.random()*Item.COLORS.length)]})
            //new THREE.MeshPhongMaterial({color: Number.parseInt(Math.floor(Math.random()*16777215).toString(16), 16)})
        );

        scene.add(this.mesh);
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

    public test(tex: THREE.Texture): void {
        this.containers[0].test(tex);
        // let target = this.group.children[0];
        // console.log(target.material);
    }
}

export default Container;