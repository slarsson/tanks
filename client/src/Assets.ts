import * as THREE from 'three';
import { MapManifest } from './Map';

interface Textures {
    warning: THREE.Texture;
    postnord: THREE.Texture;
    msc: THREE.Texture;
    maersk: THREE.Texture;
}

class Assets {
    
    public textures: Textures | undefined;
    public map: MapManifest | undefined;

    constructor(){
        this.textures = undefined;
        this.map = undefined;
    }

    public load(): Promise<boolean> {
        return new Promise(async (resolve, _) => {
            try {
                let resp = await Promise.all([
                    Assets.loadTexture('warning.png'),
                    Assets.loadTexture('pn.png'),
                    Assets.loadTexture('msc.png'),
                    Assets.loadTexture('maersk.png'),
                    Assets.loadMap('map.json')
                ]);

                this.textures = {
                    warning: resp[0],
                    postnord: resp[1],
                    msc: resp[2],
                    maersk: resp[3]
                };

                this.map = resp[4];
                
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

    private static loadMap(url): Promise<MapManifest> {
        return new Promise((resolve, reject) => {
            fetch(url)
                .then(resp => resolve(resp.json()))
                .catch(err => reject(err));
        });
    }
}

const instance = new Assets();

export {
    instance as Assets
}