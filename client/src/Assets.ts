import * as THREE from 'three';

interface Textures {
    warning: THREE.Texture;
    postnord: THREE.Texture;
    msc: THREE.Texture;
    maersk: THREE.Texture;
}

class Assets {
    
    public textures: Textures | undefined;

    constructor(){
        this.textures = undefined;
    }

    public load(): Promise<boolean> {
        return new Promise(async (resolve, _) => {
            try {
                let resp = await Promise.all([
                    Assets.loadTexture('warning.png'),
                    Assets.loadTexture('pn.png'),
                    Assets.loadTexture('msc.png'),
                    Assets.loadTexture('maersk.png'),
                ]);

                this.textures = {
                    warning: resp[0],
                    postnord: resp[1],
                    msc: resp[2],
                    maersk: resp[3]
                };
                
                resolve(true);
            } catch(err) {
                resolve(false);
            }
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

const instance = new Assets();

export {
    instance as Assets
}