// import Swag from './Swag.ts';

// let x = new Swag();

// declare let k: string;

// console.log(x.wtf());
// //console.log('hata aik');


// interface Person {
//     firstName: string;
//     lastName: string;
// }

// function greeter(person: Person) {
//     return "Hello, " + person.firstName + " " + person.lastName;
// }

// let user = { firstName: "Jane", lastName: "User" };

// document.body.textContent = greeter(user);

import Render from './Render';
import Game from './Game';
import Manager from './Manager';

import { Assets } from './AssetsTest';

// (() => {
//     const x = new Render();
// })();

window.addEventListener('DOMContentLoaded', async () => {
    await Assets.load();


    const go = new Go();

    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);

        const wasm = {
            keypress: wasm__keypress,
            poll: wasm__poll,
            update: wasm__update,
            getPos: wasm__get,
            removePlayer: wasm__removePlayer,
            setSelf: wasm__setSelf,
            //guessPosition: wasm__guessPosition,
            updateProjectiles: wasm__updateProjectiles,
            addProjectile: wasm__addProjectile
        }

        //setTimeout(() => new Manager(wasm), 3000);
        // NEW:
        new Manager(wasm);
        //new Game();
        
        
        // OLD:
        //new Render(wasm);


        (function(){var script=document.createElement('script');script.onload=function(){var stats=new Stats();document.body.appendChild(stats.dom);requestAnimationFrame(function loop(){stats.update();requestAnimationFrame(loop)});};script.src='//mrdoob.github.io/stats.js/build/stats.min.js';document.head.appendChild(script);})()
    });


});


