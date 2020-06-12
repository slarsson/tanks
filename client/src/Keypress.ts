

// type swag = {
//     [key: string]: boolean
// }

interface Keys {
    w: boolean,
    a: boolean,
    s: boolean,
    d: boolean
}

class Keypress {    
    
    status: Keys = {w: false, a: false, s: false, d: false};
    
    constructor() {
        this.add = this.add.bind(this);
        this.remove = this.remove.bind(this);

        window.addEventListener('keydown', this.add);
        window.addEventListener('keyup', this.remove);
    }

    private add(evt: KeyboardEvent): void {
        if (this.status.hasOwnProperty(evt.key)) {
            this.status[evt.key] = true;
        }
    }

    private remove(evt: KeyboardEvent): void {
        if (this.status.hasOwnProperty(evt.key)) {
            this.status[evt.key] = false;
        }
    }

    dispose(): void {

    }
}

export default Keypress;