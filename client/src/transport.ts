


// SEND:
// POLL => 
// add self to players => 0
// remove self => 1 ?
// set name => 99

import { decode } from "punycode";

// RECIVE:
//
// GAME STATE => 0
// error message => 33
// player name => 98
// kill message => 99
// remove player => 9
// set self => 10




interface Message {
    messageType: number;
}

interface ErrorMessage extends Message {
    error: string;
}

interface PlayerInfoMessage extends Message {
    id: number;
    name: string;
}

interface KillMessage extends Message {
    killerId: number;
    killedId: number;
}

interface MyIdMessage extends Message {
    id: number;
}

interface Decoder {
    myIdMessage: (arg: ArrayBuffer) => MyIdMessage; 
}





