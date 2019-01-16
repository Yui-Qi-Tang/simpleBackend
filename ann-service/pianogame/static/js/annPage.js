'use strict'
$(document).ready(function() {
    document.getElementById("btn").addEventListener("click", myFunc);
    document.getElementById("input").addEventListener("change", myFunc)
    function myFunc(){
      document.getElementById("hi").innerHTML = "Thank you. I've received your message "
    }
});