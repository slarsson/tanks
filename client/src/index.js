import Manager from './Manager';
import { Assets } from './Assets';

window.addEventListener('DOMContentLoaded', async () => {
    let data = await Promise.all([Assets.load(), fetch('main.wasm')]);
    if (data[0] != true || data[1].status != 200) {
        console.error('FAILED TO LOAD, ABORT!', data);
        return;
    }

    const go = new Go();
    let wasm = await WebAssembly.instantiateStreaming(data[1], go.importObject); 
    go.run(wasm.instance);

    new Manager({
        keypress: wasm__keypress,
        poll: wasm__poll,
        update: wasm__update,
        getPos: wasm__get,
        removePlayer: wasm__removePlayer,
        setSelf: wasm__setSelf,
        updateProjectiles: wasm__updateProjectiles,
        updateCrane: wasm__updateCrane
    });

    (function(){var script=document.createElement('script');script.onload=function(){var stats=new Stats();document.body.appendChild(stats.dom);requestAnimationFrame(function loop(){stats.update();requestAnimationFrame(loop)});};script.src='//mrdoob.github.io/stats.js/build/stats.min.js';document.head.appendChild(script);})();
});