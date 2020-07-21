import * as THREE from 'three';

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
        //"boundaries": [150, -50, 50, -50], 
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

        let xLength = Math.abs(this.manifest.boundaries[0] - this.manifest.boundaries[1]);
        let yLength = Math.abs(this.manifest.boundaries[2] - this.manifest.boundaries[3]);

        {
            
            let xCenter = this.manifest.boundaries[0] - (0.5 * xLength);
            let yCenter = this.manifest.boundaries[2] - (0.5 * yLength);

            let mesh = new THREE.PlaneGeometry(xLength, yLength, 10);
                mesh.translate(xCenter, yCenter, 0);

            this.scene.add(new THREE.Mesh(
                mesh,  
                new THREE.MeshBasicMaterial({color: 0xb7b1ae, side: THREE.DoubleSide})
            ));
        }
        
        {
            let w = 2;
            let color = 0xffff00;

            let yTop = new THREE.PlaneGeometry(xLength+2*w, w, 0);
                yTop.translate(0, this.manifest.boundaries[2] + (0.5 * w), 0);
            let yBottom = new THREE.PlaneGeometry(xLength+2*w, w, 0);
                yBottom.translate(0, this.manifest.boundaries[3] - (0.5 * w), 0);
            let xTop = new THREE.PlaneGeometry(w, yLength, 0);
                xTop.translate(this.manifest.boundaries[0] + (0.5 * w), 0, 0);
            let xBottom = new THREE.PlaneGeometry(w, yLength, 0);
                xBottom.translate(this.manifest.boundaries[1] - (0.5 * w), 0, 0);
                
            // this.scene.add(new THREE.Mesh(yTop, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide})));
            //this.scene.add(new THREE.Mesh(yBottom, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide})));
            //this.scene.add(new THREE.Mesh(xTop, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide})));
            //this.scene.add(new THREE.Mesh(xBottom, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide})));
        

            let mesh = new THREE.Mesh(yTop, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide}));
            this.scene.add(mesh);

            let mesh2 = new THREE.Mesh(xTop, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide}));
            this.scene.add(mesh2);

            let mesh3 = new THREE.Mesh(yBottom, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide}));
            this.scene.add(mesh3);

            let mesh4 = new THREE.Mesh(xBottom, new THREE.MeshBasicMaterial({color: color, side: THREE.DoubleSide}));
            this.scene.add(mesh4);

            let loader = new THREE.TextureLoader();

            // load a resource
            loader.load(
                //'postnord.png',
                'warning2.png',

                // onLoad callback
                function ( texture ) {
                    console.log('fmt:', texture.format);
                   
                    texture.wrapS = THREE.RepeatWrapping;
                    texture.wrapT = THREE.RepeatWrapping;
                    texture.minFilter = THREE.LinearFilter;
                    
                    
                    texture.repeat.set(xLength / w, 1);
                    //texture.rotation = Math.PI / 4;

                    mesh.material = new THREE.MeshBasicMaterial({
                        //color: color, 
                        //opacity: 0.5,
                        //transparent: true,
                        //side: THREE.DoubleSide,
                        map: texture
                    });

                    let t2 = texture.clone();
                    t2.needsUpdate = true;
                    console.log(t2);
                    t2.repeat.set(1, yLength / w);

                    mesh2.material = new THREE.MeshBasicMaterial({
                        //color: color, 
                        //opacity: 0.5,
                        //transparent: true,
                        side: THREE.DoubleSide,
                        map: t2
                    });

                    let t3 = texture.clone();
                    t3.needsUpdate = true;
                    //console.log(t2);
                    //t2.repeat.set(1, yLength / w);
                    mesh3.material = new THREE.MeshBasicMaterial({
                        //color: color, 
                        //opacity: 0.5,
                        //transparent: true,
                        side: THREE.DoubleSide,
                        map: t3
                    });

                    let t4 = texture.clone();
                    t4.needsUpdate = true;
                    //console.log(t2);
                    t4.repeat.set(1, yLength / w);
                    mesh4.material = new THREE.MeshBasicMaterial({
                        //color: color, 
                        //opacity: 0.5,
                        //transparent: true,
                        side: THREE.DoubleSide,
                        map: t4
                    });


                    
                    // mesh.material.color.set(0x000000);
                    // console.log(mesh.material);

                   
                    // // in this example we create the material when the texture is loaded
                    // var material = new THREE.MeshBasicMaterial( {
                    //     map: texture
                    // } );
                },

                // onProgress callback currently not supported
                undefined,

                // onError callback
                function ( err ) {
                    console.error( 'An error happened.' );
                }
            );
        }

        // {
        //     let mesh = new THREE.PlaneGeometry(1000, 1000, 10);
        //         mesh.translate(0, 0, -0.1);

        //     this.scene.add(new THREE.Mesh(
        //         mesh,
        //         new THREE.MeshBasicMaterial({color: 0x000000, side: THREE.DoubleSide})
        //     ));
        // }

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