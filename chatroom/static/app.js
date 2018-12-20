$(function () {
    var ws;
    var myData = undefined;
    var chat = $("#chat");

    if (window.WebSocket === undefined) {
        $("#container").append("Your browser does not support WebSockets");
        return;
    } else {
        ws = initWS();
    }
    function initWS() {
        var socket = new WebSocket("wss://localhost:8081/socket/handler");
        var container = $("#container");
        var chatInfo = $("#chatInfo");
        socket.onopen = function() {
            container.append("<p>Socket is open</p>");           
        };
        socket.onmessage = function (e) {
            //container.append("<p> Got :" + e.data + "</p>");
            myData = JSON.parse(e.data);
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
            }
        }
        socket.onclose = function () {
            container.append("<p>Socket closed</p>");
        }
        return socket;
    }
    $("#sendBtn").click(function (e) {
        e.preventDefault();
        let sendMsg = $("#msg").val();
        ws.send(JSON.stringify({
            //MyId: parseInt($("#numberfield").val()),
            Text: sendMsg,
            From: myData.MyId,
        }));

        chat.append(`<div class="container">
            <p>Me</p>
            <p>${sendMsg}</p>
            <span class="time-right">This time</span>
            </div>`
        );
    });

    $('#closeSocket').click(function (e) {
        e.preventDefault();
        ws.close();
    });
});