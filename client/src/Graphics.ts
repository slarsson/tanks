
const addInfoBox = (parent: HTMLElement, text: string, timeout: number | null = null): void => {
    let t: number | undefined = undefined;
    
    let root = document.createElement('div');
        root.id = 'info'; // BAD!! should be class
        root.onanimationend = () => parent.removeChild(root);
    
    let close = document.createElement('button');
        close.innerText = 'CLOSE';
        close.onclick = () => {
            clearTimeout(t);
            root.classList.add('swag')
            close.onclick = null;
        };

    let content = document.createElement('p');
        content.innerText = text;
    
    if (timeout != null) {
        t = window.setTimeout(() => {
            root.classList.add('swag')
            close.onclick = null;
        }, timeout);
    }

    root.appendChild(content);
    root.appendChild(close);
    parent.appendChild(root);
};

const addKillMessage = (parent: HTMLElement, killer: string, killed: string): void => {
    let root = document.createElement('div');
        root.innerText = killer + ' KILLED ' + killed;

    setTimeout(() => parent.removeChild(root), 2000);
};


class Graphics {

    public static readonly INFO = 0;
    public static readonly WARNING = 1;
    public static readonly ERROR = 2;

    private root: HTMLElement;

    private info: HTMLElement;
    private killLog: HTMLElement;

    constructor() {
        let rootDiv = document.getElementById('graphics');
        if (rootDiv != null) {
            this.root = rootDiv;
        } else {
            console.error('GRAPHICS: root div not found');
            this.root = document.createElement('div');
            this.root.id = 'graphics';
        }

        this.info = document.createElement('div');
        this.info.id = 'info';

        this.killLog = document.createElement('div');
        this.killLog.id = 'kills';

        this.root.appendChild(this.info);
        this.root.appendChild(this.killLog);

        //addInfoBox(this.root, 'hata aik');
        
        // let x = new InfoBox(this.root);
        // x.add('wtf?');

        // this.root.classList.add('blur');


        // console.log('NEW GRAPHICS CREATED...');

        let swag = document.createElement('div');
        swag.innerHTML = `
        <div class="center">
            <div class="panel">
                <div class="logo">
                    game_my_game
                </div>


                <form>
                    <div class="add-name">
                        <label for="">Enter player name</label>
                        <input type="text"/>
                        <button type="submit">continue</button>
                    </div>
                </form>
            </div>
        </div>
        `;

        this.root.appendChild(swag);
        //this.root.classList.add('blur');
    }

    newInfoBox(text: string, timeout: number | null = null, type: number = 0): void {
        let t: number | undefined = undefined;
        
        let icon = document.createElement('div');
            icon.classList.add('icon');

        let root = document.createElement('div');
            if (type == Graphics.INFO) {
                root.classList.add('info');
                icon.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"/></svg>';
            } else if (type == Graphics.WARNING) {
                icon.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/></svg>';
                root.classList.add('warning');
            } else if (type == Graphics.ERROR) {
                icon.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/></svg>';
                root.classList.add('error');
            }

            root.classList.add('info-box', 'fade-in'); // BAD!! should be class
            root.onanimationend = (e: AnimationEvent) => {
                if (e.animationName == 'out') {
                    this.info.removeChild(root);
                }
            }
            
            //root.onanimationend = () => this.info.removeChild(root);
        
        let close = document.createElement('button');
            close.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>';
            close.onclick = () => {
                clearTimeout(t);
                //root.onanimationend = () => this.info.removeChild(root);
                root.classList.remove('fade-in');
                root.classList.add('fade-out');
                close.onclick = null;
            };

        let content = document.createElement('p');
            content.innerText = text;
        
        
            //icon.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/></svg>';
            //icon.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/></svg>';

        if (timeout != null) {
            t = window.setTimeout(() => {
                root.classList.remove('fade-in');
                root.classList.add('fade-out');
                close.onclick = null;
            }, timeout);
        }

        root.appendChild(icon);
        root.appendChild(content);
        root.appendChild(close);
        this.info.appendChild(root);
    }

    newKillMessage(killer: string, killed: string, timeout: number = 2000): void {
        let root = document.createElement('div');
        
        let k1 = document.createElement('p');
            k1.innerText = killer;
        let k2 = document.createElement('p');
            k2.innerText = killed;

        let div = document.createElement('div');
            div.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" width="70" height="20" viewBox="0 0 70 20" > <rect x="0" y="10" width="60" height="10"/> <rect x="20" y="0" width="20" height="10"/><rect x="40" y="3" width="30" height="4"/></svg>';

        
        root.appendChild(k1);
        root.appendChild(div);
        root.appendChild(k2);
        this.killLog.appendChild(root);

        setTimeout(() => this.killLog.removeChild(root), timeout);
    }


}

export default Graphics;