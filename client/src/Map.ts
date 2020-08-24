import * as THREE from 'three';
import { Assets } from './Assets';
import { Container } from './Container';
import Crane from './Crane';

interface ContainersTest {
    position: {x: number, y: number, z: number};
    rotation: number;
    total: number;
    bottom: number;
} 

interface MapManifest {
    name: string;
    boundaries: number[];
    containers: ContainersTest[];
}

class GameMap {

    private manifest: MapManifest;
    private scene: THREE.Scene;
    private crane: Crane;

    constructor(scene: THREE.Scene) {
        this.scene = scene;
            
        if (Assets.map != undefined) {
            this.manifest = Assets.map; // clone?
        } else {
            this.manifest = {
                name: 'default',
                boundaries: [12, -12, 12, -12],
                containers: []
            };
        }

        for (const item of this.manifest.containers) {
            let c = new Container(this.scene, item.total, item.bottom);
                c.setPosition(item.position.x, item.position.y, item.position.z);
                c.setRotation(item.rotation);
        }

        let xLength = Math.abs(this.manifest.boundaries[0] - this.manifest.boundaries[1]);
        let yLength = Math.abs(this.manifest.boundaries[2] - this.manifest.boundaries[3]);

        {
            
            let xCenter = this.manifest.boundaries[0] - (0.5 * xLength);
            let yCenter = this.manifest.boundaries[2] - (0.5 * yLength);

            let mesh = new THREE.PlaneGeometry(xLength, yLength, 10);
                mesh.translate(xCenter, yCenter, 0);

            this.scene.add(new THREE.Mesh(
                mesh,
                //new THREE.MeshBasicMaterial({color: 0xc6c1b0, side: THREE.DoubleSide})  
                //new THREE.MeshBasicMaterial({color: 0xa0a09a, side: THREE.DoubleSide})  
                new THREE.MeshPhongMaterial({color: 0xb7b1ae, side: THREE.DoubleSide})
            ));
        }
        
        {
            let w = 2;

            let yTop = new THREE.PlaneGeometry(xLength+2*w, w, 0);
                yTop.translate(0, this.manifest.boundaries[2] + (0.5 * w), 0);
            let yBottom = new THREE.PlaneGeometry(xLength+2*w, w, 0);
                yBottom.translate(0, this.manifest.boundaries[3] - (0.5 * w), 0);
            let xTop = new THREE.PlaneGeometry(w, yLength, 0);
                xTop.translate(this.manifest.boundaries[0] + (0.5 * w), 0, 0);
            let xBottom = new THREE.PlaneGeometry(w, yLength, 0);
                xBottom.translate(this.manifest.boundaries[1] - (0.5 * w), 0, 0);
            
            if (Assets.textures?.warning != undefined) {
                let yTexture = Assets.textures?.warning.clone();
                    yTexture.wrapS = THREE.RepeatWrapping;
                    yTexture.wrapT = THREE.RepeatWrapping;
                    yTexture.repeat.set((xLength / w), 1);
                    yTexture.needsUpdate = true;
                    yTexture.minFilter = THREE.LinearFilter;

                let xTexture = Assets.textures?.warning.clone();
                    xTexture.wrapS = THREE.RepeatWrapping;
                    xTexture.wrapT = THREE.RepeatWrapping;
                    xTexture.repeat.set(1, (yLength / w));    
                    xTexture.needsUpdate = true;
                    xTexture.minFilter = THREE.LinearFilter;

                this.scene.add(new THREE.Mesh(yTop, new THREE.MeshBasicMaterial({map: yTexture})));
                this.scene.add(new THREE.Mesh(yBottom, new THREE.MeshBasicMaterial({map: yTexture})));
                this.scene.add(new THREE.Mesh(xTop, new THREE.MeshBasicMaterial({map: xTexture})));
                this.scene.add(new THREE.Mesh(xBottom, new THREE.MeshBasicMaterial({map: xTexture})));
            }
        }

        {
            let hemiLight = new THREE.HemisphereLight(0xffffff, 0xffffff, 0.7);
            hemiLight.position.set(0, 500, 0);
            this.scene.add( hemiLight );

            let dirLight = new THREE.DirectionalLight(0xffffff, 0.6);
            dirLight.position.set(-1, 0.75, 1);
            dirLight.position.multiplyScalar(5);
            this.scene.add(dirLight);
        }

        this.crane = new Crane(this.scene);

//         console.log(Assets.objects?.crane);
// //         if (Assets.objects != undefined) {
// //             const crane = Assets.objects.crane.scene;
// //             //crane.rotation.y = Math.random() * 3;
// //             crane.rotation.x = Math.PI / 2;

// //             crane.scale.set(0.5, 0.5, 0.5);
// //             crane.position.x = 20;
            
// //             let child = crane.children[0];
// //             if (child instanceof THREE.Mesh) {
// //                 child.material = new THREE.MeshPhongMaterial({color: 0xff0000, side: THREE.DoubleSide}); 
// //             }

// // //            console.log(crane.children[0].isMesh);
            
            
            
//             //this.scene.add(crane);
//             //this.scene.add(Assets.objects.crane.scene);
//         }

        //scene.add(new THREE.AxesHelper(150)); 
    }

    outOfMap(x: number, y: number): boolean {
        return x > this.manifest.boundaries[0] || x < this.manifest.boundaries[1] || y > this.manifest.boundaries[2] || y < this.manifest.boundaries[3];
    }

    update(dt: number) {
        this.crane.update(dt);
    }
}

export {
    GameMap,
    MapManifest
};