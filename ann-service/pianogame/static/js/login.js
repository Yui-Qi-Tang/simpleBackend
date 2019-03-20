'use strict'
let sendLogin = (url, account, pwd) => {
    // set data
    let userData = {
        "account": account,
        "password": pwd
    };
    const fetchSettings = {
        method: "POST",
        headers: {
            'Content-Type': 'application/json; charset=UTF-8'
        },
        body: JSON.stringify(userData)
    };

    // fetch
    fetch(url, fetchSettings)
    .then((response) => {
        window.location.replace("/game");
        return response.json(); 
    })
    .then((jsonData) => {
        alert(jsonData.msg);
    })
    .catch((err) => {
        console.log(err);
    });
};

let loginEvent = () => {
    const loginURL = "/login";
    const account = document.getElementById("account").value;
    const password = document.getElementById("password").value;
    sendLogin(loginURL, account, password);

};
// wait for html ready!
document.addEventListener('DOMContentLoaded', function() { 
    let login = document.getElementById("login");
    login.addEventListener("click", loginEvent);
}, false);