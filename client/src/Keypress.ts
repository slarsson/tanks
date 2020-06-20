

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
        this.poll = this.poll.bind(this);

        // window.addEventListener('keydown', this.add);
        // window.addEventListener('keyup', this.remove);
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

    poll(): Uint8Array {
        let x = new Uint8Array(4);
        x[0] = (this.status.w ? 1 : 0);
        x[1] = (this.status.a ? 1 : 0);
        x[2] = (this.status.s ? 1 : 0);
        x[3] = (this.status.d ? 1 : 0);
        return x;
        // return (
        //     (this.status.w ? 1 : 0).toString() + 
        //     (this.status.a ? 1 : 0).toString() + 
        //     (this.status.s ? 1 : 0).toString() + 
        //     (this.status.d ? 1 : 0).toString()
        // );
    }

    dispose(): void {

    }
}

export default Keypress;