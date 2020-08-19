import * as THREE from 'three';

import { Assets } from './AssetsTest';
import Container from './Container';
import { Mesh } from 'three';

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
                    "rotation": 0,
                    "total": 25,
                    "bottom": 10
                },
                {
                    "position": {
                        "x": -20,
                        "y": 15,
                        "z": 0
                    },
                    "rotation": -0.3,
                    "total": 5,
                    "bottom": 5
                },
                {
                    "position": {
                        "x": 0,
                        "y": 30,
                        "z": 0
                    },
                    "rotation": 0.5,
                    "total": 2,
                    "bottom": 1
                },
                {
                    "position": {
                        "x": -35,
                        "y": -30,
                        "z": 0
                    },
                    "rotation": 1.2,
                    "total": 9,
                    "bottom": 5
                }
            ]
        }`);

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
                new THREE.MeshBasicMaterial({color: 0xa0a09a, side: THREE.DoubleSide})  
                //new THREE.MeshPhongMaterial({color: 0xb7b1ae, side: THREE.DoubleSide})
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
            let hemiLight = new THREE.HemisphereLight(0xffffff, 0xffffff, 0.7);
            hemiLight.position.set(0, 500, 0);
            this.scene.add( hemiLight );

            let dirLight = new THREE.DirectionalLight(0xffffff, 0.6);
            dirLight.position.set(-1, 0.75, 1);
            dirLight.position.multiplyScalar(5);
            this.scene.add(dirLight);
        }

        //scene.add(new THREE.AxesHelper(150));   
        
        // {
        //     let geometry = new THREE.Geometry();
        //         geometry.vertices= [
        //             // new THREE.Vector3(2,1,0), 
        //             // new THREE.Vector3(1,3,0), 
        //             // new THREE.Vector3(3,4,0)
        //             new THREE.Vector3(0,0,0), 
        //             new THREE.Vector3(0.5,0,1), 
        //             new THREE.Vector3(-0.5,0,1)
        //         ]; 
        //         geometry.faces = [new THREE.Face3(0, 1, 2)];

        //     let swag = new THREE.Mesh(
        //         geometry,
        //         new THREE.MeshBasicMaterial({color: 0x00ff00, transparent: true, opacity: 0.5, side: THREE.DoubleSide})
        //     );

        //     swag.position.z = 4;
        //     swag.scale.set(2, 2, 2);

        //     this.scene.add(swag);
        // }
    }

    outOfMap(x: number, y: number): boolean {
        return x > this.manifest.boundaries[0] || x < this.manifest.boundaries[1] || y > this.manifest.boundaries[2] || y < this.manifest.boundaries[3];
    }
}

export default GameMap;