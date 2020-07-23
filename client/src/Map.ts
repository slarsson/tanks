import * as THREE from 'three';

import { Assets } from './AssetsTest';
import Container from './Container';

interface ContainersTest {
    position: {x: number, y: number, z: number};
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

    constructor(scene: THREE.Scene) {
        this.scene = scene;
        this.manifest = JSON.parse(`{
            "name": "Port of Nrkp",
            "boundaries": [50, -50, 50, -50],
            "containers": [
                {
                    "position": {
                        "x": 40,
                        "y": 0,
                        "z": 0
                    },
                    "total": 25,
                    "bottom": 10
                },
                {
                    "position": {
                        "x": -20,
                        "y": 15,
                        "z": 0
                    },
                    "total": 5,
                    "bottom": 5
                },
                {
                    "position": {
                        "x": 0,
                        "y": 30,
                        "z": 0
                    },
                    "total": 2,
                    "bottom": 1
                }
            ]
        }`);

        // let c = new Container(this.scene, 1, 1);
        // c.setPosition(0, 0, 5);
        for (const item of this.manifest.containers) {
            let c = new Container(this.scene, item.total, item.bottom);
                c.setPosition(item.position.x, item.position.y, item.position.z);
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
                //new THREE.MeshBasicMaterial({color: 0xa0a09a, side: THREE.DoubleSide})  
                new THREE.MeshBasicMaterial({color: 0xb7b1ae, side: THREE.DoubleSide})
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
                let xTexture = Assets.textures?.warning.clone();
                    xTexture.wrapS = THREE.RepeatWrapping;
                    xTexture.wrapT = THREE.RepeatWrapping;
                    xTexture.repeat.set(1, (yLength / w));    
                    xTexture.needsUpdate = true;
                
                this.scene.add(new THREE.Mesh(yTop, new THREE.MeshBasicMaterial({map: yTexture})));
                this.scene.add(new THREE.Mesh(yBottom, new THREE.MeshBasicMaterial({map: yTexture})));
                this.scene.add(new THREE.Mesh(xTop, new THREE.MeshBasicMaterial({map: xTexture})));
                this.scene.add(new THREE.Mesh(xBottom, new THREE.MeshBasicMaterial({map: xTexture})));
            }
        }

        {
            let light = new THREE.AmbientLight(0x404040, 2); // soft white light
            this.scene.add( light );

            let hemiLight = new THREE.HemisphereLight(0x5780EA, 0xffffff, 0.3);
            this.scene.add(hemiLight);

            // let light = new THREE.SpotLight(0xffa95c,2);
            // light.position.set(-50,50,50);
            // light.castShadow = true;
            // this.scene.add( light );

            // let l1 = new THREE.PointLight(0xff0000, 0.5, 100);
            // l1.position.set(10, 10, 0);
            // this.scene.add(l1);

            // let l2 = new THREE.DirectionalLight(0xefefff, 1.4);
            // l2.position.set(1, 1, 1).normalize();
            // l2.castShadow = true;
            // this.scene.add(l2);
        }
    }
}

export default GameMap;