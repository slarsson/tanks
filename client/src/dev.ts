import * as THREE from 'three';

const helper = (scene: THREE.Scene) => {
    scene.add(new THREE.AxesHelper(15));

    let l1 = new THREE.PointLight(0xff0000, 0.5, 100);
    l1.position.set(10, 10, 0);
    scene.add(l1);

    let l2 = new THREE.DirectionalLight(0xefefff, 1.4);
    l2.position.set(1, 1, 1).normalize();
    l2.castShadow = true;
    scene.add(l2);

    scene.add(new THREE.Mesh(
        new THREE.CylinderGeometry(1700, 1700, 2000, 50),  
        new THREE.MeshBasicMaterial({color: 0x5780EA, side: THREE.BackSide})
    ));

    let plane = new THREE.Mesh(
        new THREE.PlaneGeometry(100, 100, 10),  
        new THREE.MeshBasicMaterial( {color: 0xdf9e6f, side: THREE.DoubleSide} )
    )
    plane.position.z = -0.5;
    scene.add(plane);
};

const obstacleTest = (scene: THREE.Scene) => {
    let jsonInput = `{
        "name": "Biltema",
		"blocks": [
			{
				"name": "wall",
				"coords": [[0, 16, 0], [10, 16, 0], [10, 15, 0], [0, 15, 0]]
			},
			{
				"name": "house1",
				"coords": [[10, 10, 0], [20, 10, 0], [20, 0, 0], [10, 0, 0]]
			}
		]
    }`;
    
    let wtf = JSON.parse(jsonInput);
    console.log(JSON.parse(jsonInput));

    for (const item of wtf.blocks) {
        
        // for (const lol of item.coords) {
        //     console.log(lol);
        // }
    }

    
    let wall = new THREE.Mesh(
        new THREE.BoxGeometry(10, 1, 5),
        new THREE.MeshPhongMaterial({color: 0x4d5858})
    );
    wall.position.set(5, 15.5, 2);
    scene.add(wall);

    let house = new THREE.Mesh(
        new THREE.BoxGeometry(10, 10, 5),
        new THREE.MeshPhongMaterial({color: 0x4d5858})
    );
    house.position.set(15, 5, 2);
    scene.add(house);
};

export {
    helper,
    obstacleTest
}