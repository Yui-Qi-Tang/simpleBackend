'use strict'

let support_websocket = () => {
    if (window.WebSocket === undefined) {
        alert("Your browser does not support WebSockets");
        return false;
    }
    return true;
};

let socketClientEvents = {
    open: () => {
        console.log("Socket open!!");
    },
    receiveMsg: function(e) {
        let myData = JSON.parse(e.data);
        console.log(myData);
        if(myData.Key) {
            playPianoKey(myData.Key);
        }
        /*
        if(!myData.From) {
          chatInfo.append(`<h2>${myData.Text}</h2>`);
          chatInfo.append(`This is your id: <h2>${myData.MyId}</h2>`);
        }
        else{
          chat.append(
            `<div class="container darker">
             <p>${myData.From} says: </p>
             <p>${myData.Text}</p>
             <span class="time-left">11:01</span>
             </div>`
          );  
        }*/
    },
    destroy: () => {
        console.log("close socket");
    },
    error: () => {
        throw 'socket connent error';
    },
}

/**
 * TODO: create socket and return
 */
class GameWebSocketFactory {
    constructor(socketServerURL) {
        this.serverURL = socketServerURL;
    }

    create() {
        if(this.serverURL == undefined || this.serverURL == null) {
            return false;
        }
        let socket = new WebSocket(this.serverURL);
        socket.onopen = socketClientEvents.open;
        socket.onmessage = socketClientEvents.receiveMsg;
        socket.onclose = socketClientEvents.destroy;
        socket.onerror = socketClientEvents.error;
        return socket;
    }
}

function start_connect_to_server() {
    return new Promise(function(resolve, reject) {
        try {
            let s = new GameWebSocket("wss://127.0.0.1:8081/game/socket");
            resolve(s);
        } catch(e) {
            reject(e);
        }
    });
}