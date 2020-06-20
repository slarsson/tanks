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

// (() => {
//     const x = new Render();
// })();

window.addEventListener('DOMContentLoaded', () => {
    
    const go = new Go();

    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
        go.run(result.instance);

        const wasm = {
            test: swag,
            print: wasmprint,
            state: state,
            poll: poll,
            update: wasmupdate,
            getPos: wasmgetpos,
            local: wasmglocal
        }

        wasm.test();
        new Render(wasm);
    });


});


