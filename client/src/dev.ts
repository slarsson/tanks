import * as THREE from 'three';

import Container from './Container';

const helper = (scene: THREE.Scene) => {
    scene.add(new THREE.AxesHelper(15));

    // let l1 = new THREE.PointLight(0xff0000, 0.5, 100);
    // l1.position.set(10, 10, 0);
    // scene.add(l1);

    // let l2 = new THREE.DirectionalLight(0xefefff, 1.4);
    // l2.position.set(1, 1, 1).normalize();
    // l2.castShadow = true;
    // scene.add(l2);

    // scene.add(new THREE.Mesh(
    //     new THREE.CylinderGeometry(1700, 1700, 2000, 50),  
    //     new THREE.MeshBasicMaterial({color: 0x5780EA, side: THREE.BackSide})
    // ));

    // let plane = new THREE.Mesh(
    //     new THREE.PlaneGeometry(100, 100, 10),  
    //     new THREE.MeshBasicMaterial( {color: 0xdf9e6f, side: THREE.DoubleSide} )
    // )
    // plane.position.z = -0.5;
    // scene.add(plane);

    // let plane2 = new THREE.Mesh(
    //     new THREE.PlaneGeometry(140, 140, 10),  
    //     new THREE.MeshBasicMaterial( {color: 0xff0000, side: THREE.DoubleSide} )
    // )
    // plane2.position.z = -0.6;
    // scene.add(plane2);
};

const obstacleTest = (scene: THREE.Scene) => {
    // let jsonInput = `{
    //     "name": "Biltema",
	// 	"blocks": [
	// 		{
	// 			"name": "wall",
	// 			"coords": [[0, 16, 0], [10, 16, 0], [10, 15, 0], [0, 15, 0]]
	// 		},
	// 		{
	// 			"name": "house1",
	// 			"coords": [[10, 10, 0], [20, 10, 0], [20, 0, 0], [10, 0, 0]]
	// 		}
	// 	]
    // }`;
    
    // let wtf = JSON.parse(jsonInput);
    // console.log(JSON.parse(jsonInput));

    // for (const item of wtf.blocks) {
        
    //     // for (const lol of item.coords) {
    //     //     console.log(lol);
    //     // }
    // }

    
    // let wall = new THREE.Mesh(
    //     new THREE.BoxGeometry(10, 1, 5),
    //     new THREE.MeshPhongMaterial({color: 0x4d5858})
    // );
    // wall.position.set(5, 15.5, 2);
    // scene.add(wall);

    let c1 = new Container(scene, 25, 10);
        c1.setPosition(40, 0, 0);
    let c2 = new Container(scene, 5, 5);
        c2.setPosition(-20, 15, 0);
    // let c3 = new Container(scene, 2, 1);
    //     c3.setPosition(0, 30, 0);


        let loader = new THREE.TextureLoader();

        loader.load(
            'postnord.png',
            //'warning.jpeg',

            // onLoad callback
            function ( texture ) {
                //console.log(texture);
                let c3 = new Container(scene, 2, 1);
                c3.setPosition(0, 30, 0);
                c3.test(texture);



               
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


    // let house = new THREE.Mesh(
    //     new THREE.BoxGeometry(10, 10, 5),
    //     new THREE.MeshPhongMaterial({color: 0x4d5858})
    // );
    // house.position.set(15, 5, 2);
    // scene.add(house);


    // // wtf:
    // let geometry = new THREE.BufferGeometry();
    // // create a simple square shape. We duplicate the top left and bottom right
    // // vertices because each vertex needs to appear once per triangle.
    // let vertices = new Float32Array( [
    //     0, 16, 1, // first
    //     10, 16, 1, 
    //     10, 15, 1,
        
    //     0, 16, 1, // first
    //     10, 15, 1, // last of prev verticies
    //     0, 15, 1,

        
    //     // -10.0, -1.0,  1.0,
    //     // 1.0, -1.0,  1.0,
    //     // 1.0,  5.0,  1.0,

    //     // 1.0,  1.0,  1.0,
    //     // -1.0,  1.0,  1.0,
    //     // -1.0, -1.0,  1.0
    // ] );

    // // itemSize = 3 because there are 3 values (components) per vertex
    // geometry.setAttribute( 'position', new THREE.BufferAttribute( vertices, 3 ) );
    // let material = new THREE.MeshBasicMaterial( { color: 0xff0000, side: THREE.DoubleSide } );
    // let mesh = new THREE.Mesh( geometry, material );
    // scene.add(mesh);
};

export {
    helper,
    obstacleTest
}