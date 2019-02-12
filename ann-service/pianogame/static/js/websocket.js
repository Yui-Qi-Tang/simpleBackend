'use strict'

let support_websocket = () => {
    if (window.WebSocket === undefined) {
        alert("Your browser does not support WebSockets");
        return false;
    }
    return true;
};


class GameWebSocket {
    constructor(socketServerURL) {
        this.socket = new WebSocket(socketServerURL);
    }
}
