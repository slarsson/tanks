import * as THREE from 'three';

interface Textures {
    warning: THREE.Texture;
    postnord: THREE.Texture;
}

class AssetsTest {
    
    //private data: AssetsData | undefined;

    public textures: Textures | undefined;

    constructor(){
        this.textures = undefined;
    }

    public load(): Promise<null> {
        return new Promise(async (resolve, reject) => {
            let test = await AssetsTest.loadTexture('warning2.png');            
            //let pn = await AssetsTest.loadTexture('warning.jpeg');            

            let pn = await AssetsTest.loadTexture('pn.png');  

            this.textures = {
                warning: test,
                postnord: pn
            }

            console.log(this.textures);
            resolve();
        });
    }

    private static loadTexture(url): Promise<THREE.Texture> {
        return new Promise((resolve, reject) => {
            new THREE.TextureLoader().load(
                url,
                (texture) => resolve(texture),
                undefined,
                (err) => reject(err)
            );
        });
    }
}

const instance = new AssetsTest();
//Object.freeze(instance);

//export default instance;

export {
    instance as Assets
}
  