interface NameInputActions {
    showError: (arg: string) => void;
    hideError: () => void;
    setLoading: (arg: boolean) => void;
} 

class NameInput {

    private root: HTMLElement;
    private status: HTMLElement;
    private error: HTMLElement;
    private input: HTMLInputElement;
    private loading: boolean;
    private callback: (arg: string) => void;

    constructor(_root, cb) {
        this.root = document.createElement('div');
        _root.appendChild(this.root);

        this.callback = cb;
        this.status = document.createElement('div');
        this.error = document.createElement('div');
        this.loading = false;
        
        this.submit = this.submit.bind(this);
        this.setLoading = this.setLoading.bind(this);

        // input 
        this.input = document.createElement('input');
        this.input.id = '_player';
        this.input.placeholder = 'krillex';
        this.input.type = 'text';
        this.input.oninput = () => {
            this.input.classList.remove('input-error');
            this.error.innerHTML = '';
        }

        // label
        let label = document.createElement('label');
            label.innerText = 'Enter player name';
            label.htmlFor = this.input.id;

        // form
        let form = document.createElement('form');
            form.autocomplete = 'off';
            form.onsubmit = this.submit;

        // logo
        let logo = document.createElement('div');
            logo.classList.add('logo');

        // containerz
        let container = document.createElement('div');
            container.classList.add('popup--container');

        let window = document.createElement('div');
            window.classList.add('popup--window');

        let innerContainer = document.createElement('div');
            innerContainer.classList.add('window--container');

        container.appendChild(window);
        window.appendChild(logo);
        window.appendChild(form);
        form.appendChild(innerContainer);
        innerContainer.appendChild(label);
        innerContainer.appendChild(this.input);
        innerContainer.appendChild(this.error);
        innerContainer.appendChild(this.status);
        this.status.appendChild(NameInput.createSubmitButton());

        this.root.appendChild(container);
    }


    private static createSubmitButton(loading: boolean = false): HTMLElement {
        if (!loading) {
            let button = document.createElement('button');
                button.type = 'submit';
                button.classList.add('submit-button');
                button.innerHTML = `
                    <p>PLAY GAME</p>
                    <svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M8 5v14l11-7z"/></svg>
                `;
            return button;
        }

        let div = document.createElement('div');
            div.classList.add('submit-button', 'submit-button-loading');
            div.innerHTML = `<div class="loading"></div>`;
        return div;
    }

    private submit(evt: Event): void {
        evt.preventDefault();
        if (this.loading) {return;}
        
        let value = this.input.value;

        if (!/^[A-Za-z0-9_-]+$/.test(value)) {
            this.showError('Letters and numbers only');
            return;
        }

        if (value.length > 20) {
            this.showError(`Name to long, ${value.length} char. [max: 20]`);
            return;
        }

        this.setLoading(true);
        this.callback(this.input.value);  
    }

    setLoading(state: boolean): void {
        this.loading = state;
        this.status.innerHTML = '';
        this.status.appendChild(NameInput.createSubmitButton(state));
    }

    showError(msg: string): void {
        let error = document.createElement('div');
            error.classList.add('error');
            
        let errorText = document.createElement('p');
            errorText.innerText = msg;
        
        error.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/></svg>';
        error.appendChild(errorText);
        
        this.input.classList.remove('input-error');
        this.input.offsetHeight; // force repaint, to recognize that animation-state has changed
        this.input.classList.toggle('input-error');

        this.error.innerHTML = '';
        this.error.appendChild(error);
        this.setLoading(false);
    }

    hideError(): void {
        this.error.innerHTML = '';
    }

    dispose(): void {
        this.root.remove();
    }
}

class Graphics {

    public static readonly INFO = 0;
    public static readonly WARNING = 1;
    public static readonly ERROR = 2;

    private root: HTMLElement;
    private info: HTMLElement;
    private killLog: HTMLElement;
    private nameInput: NameInput | null;
    private connectedPlayers: HTMLElement;
    private outOfMap: HTMLElement | null = null;
    private outOfMapTimeout: number = -1;

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
        this.info.classList.add('info--container');

        this.killLog = document.createElement('div');
        this.killLog.classList.add('kill--container');

        this.root.appendChild(this.info);
        this.root.appendChild(this.killLog);

        this.nameInput = null;

        this.connectedPlayers = document.createElement('div');
        this.connectedPlayers.classList.add('players');
        this.connectedPlayers.title = 'connected players';
        this.root.appendChild(this.connectedPlayers);
    }

    newNameInput(cb: (arg: string) => void): NameInputActions {
        this.nameInput = new NameInput(this.root, cb);
        return {
            showError: (arg: string) => {this.nameInput != null ? this.nameInput.showError(arg) : null},
            hideError: () => {this.nameInput != null ? this.nameInput.hideError() : null},
            setLoading: (arg: boolean) => {this.nameInput != null ? this.nameInput.setLoading(arg) : null},
        };
    }

    removeNameInput(): void {
        this.nameInput?.dispose();
        this.nameInput = null;
    }

    addMessage(text: string, timeout: number | null = null, type: number = 0): void {
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

            root.classList.add('info-box', 'fade-in');
            root.onanimationend = (e: AnimationEvent) => {
                if (e.animationName == 'out') {
                    this.info.removeChild(root);
                }
            }
        
        let close = document.createElement('button');
            close.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 0 24 24" width="24"><path d="M0 0h24v24H0z" fill="none"/><path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>';
            close.onclick = () => {
                clearTimeout(t);
                root.classList.remove('fade-in');
                root.classList.add('fade-out');
                close.onclick = null;
            };

        let content = document.createElement('p');
            content.innerText = text;
        
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

    addKillMessage(killer: string, killed: string, timeout: number = 2000): void {
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
    
    setConnectedPlayers(list: string[]): void {
        this.connectedPlayers.innerText = `${list.join(', ')} (${list.length} connected)`; 
    }

    showOutOfMap(n: number): void {
        this.outOfMap = document.createElement('div');
        this.outOfMap.classList.add('outofmap');
        let box = document.createElement('div');
            box.innerText = 'RETURN TO MAP';
        let counter = document.createElement('div');
            counter.innerText = '5';

        box.appendChild(counter);
        this.outOfMap.appendChild(box);
        this.root.appendChild(this.outOfMap);

        let countdown = (n: number) => {
            counter.innerText = n.toString();
            n--;
            if (n == -1) return;
            this.outOfMapTimeout = window.setTimeout(() => countdown(n), 1000);
        };

        countdown(n);
    }

    hideOutOfMap(): void {
        window.clearTimeout(this.outOfMapTimeout);
        if (this.outOfMap != null) {
            this.root.removeChild(this.outOfMap);
            this.outOfMap = null;
        }
    }
}

export default Graphics;