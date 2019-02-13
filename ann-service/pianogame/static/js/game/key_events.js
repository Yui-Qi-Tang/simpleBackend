'use strict'

var fired = false; // global to denote status of key up/down?

function Define(event){
    if(event.keyCode == 65 || event.keyCode == 87 || event.keyCode == 83 || event.keyCode == 69 || event.keyCode == 68 || event.keyCode == 70 || event.keyCode == 84 || event.keyCode == 71 || event.keyCode == 89 || event.keyCode == 72 || event.keyCode == 85 || event.keyCode == 74 || event.keyCode == 75 || event.keyCode == 79 || event.keyCode == 76|| event.keyCode == 80 || event.keyCode == 186){
        // SetMeasure2();
        playNote(event);
    } else {
        return;
    }
};

function playNote(event) {
    if(!fired){//這裡的 (!fired) = true，所以會繼續執行 if 後面的 Code，但在 Code 的第一行已經把 fired 變成 true，所以下次的 (!fired) = false，就會跳出 if 的執行圈了。  
        fired = true;
        var n1 = event.timeStamp;
        console.log(n1);
        document.getElementById('showtime').innerHTML = n1;
    };
    const audio = document.querySelector(`audio[data-key="${event.keyCode}"]`),
            key = document.querySelector(`.key[data-key="${event.keyCode}"]`);
    if (!key) {
        return;
    }
    const keyNote = key.getAttribute("data-note");
    key.classList.add("playing");
    note.innerHTML = keyNote;
    audio.currentTime = 0;
    audio.play();
    
};

function playPianoKey(keyNum) {
    const audio = document.querySelector(`audio[data-key="${keyNum}"]`);
    const key = document.querySelector(`.key[data-key="${keyNum}"]`);
    if (!key) {
        return;
    } // fi
    const keyNote = key.getAttribute("data-note");
    key.classList.add("playing");
    note.innerHTML = keyNote;
    audio.currentTime = 0;
    audio.play();
}

function Keyup(event){
    if(event.keyCode == 65 || event.keyCode == 87 || event.keyCode == 83 || event.keyCode == 69 || event.keyCode == 68 || event.keyCode == 70 || event.keyCode == 84 || event.keyCode == 71 || event.keyCode == 89 || event.keyCode == 72 || event.keyCode == 85 || event.keyCode == 74 || event.keyCode == 75 || event.keyCode == 79 || event.keyCode == 76|| event.keyCode == 80 || event.keyCode == 186){
        KeyupCountTime(event);
    } else {
        return;
    };
    function KeyupCountTime(event){
        var n2 = event.timeStamp;
        //console.log(n2);
        var n1 = document.getElementById('showtime').innerHTML;
        var timelapse =((n2-n1)/1000);
        var rate = document.getElementById('MusicRate').value;
        var sheetnumber = document.getElementById('sheetnumber').textContent; // I think this is 五線譜？
        // console.log(timelapse); 
        //Rate
        var beat2 = (60/rate)*2;
        var beat4 = (60/rate);
        var beat8 = (60/rate/2);
        var beat15 = (60/rate/4);

        var keynumber2;
        switch (event.key){
          case 'a' :
            keynumber2 = 0;
            break;
          case 'w' :
            keynumber2 = 1;
            break;
          case 's' :
            keynumber2 = 2;
            break;
          case 'e' :
            keynumber2 = 3;
            break;
          case 'd' :
            keynumber2 = 4;
            break;
          case 'f' :
            keynumber2 = 5;
            break;
          case 't' :
            keynumber2 = 6;
            break;
          case 'g' :
            keynumber2 = 7;
            break;
          case 'y' :
            keynumber2 = 8;
            break;
          case 'h' :
            keynumber2 = 9;
            break;
          case 'u' :
            keynumber2 = 10;
            break;
          case 'j' :
            keynumber2 = 11;
            break;
          case 'k' :
            keynumber2 = 12;
            break;
          case 'o' :
            keynumber2 = 13;
            break;
          case 'l' :
            keynumber2 = 14;
            break;
          case 'p' :
            keynumber2 = 15;
            break;
          case ';' :
            keynumber2 = 16;
            break; 
        };
/*      Buggy so comment out
        if(timelapse > beat2){
            document.getElementById(`sheet${sheetnumber}`).innerHTML += `<img src ='https://img.icons8.com/windows/32/000000/musical.png' style='width:50px; height:30px;position:relative;margin-left:3%;top:${(52-keynumber2)}%;'>`
        }else if(timelapse > beat4 && timelapse < beat2){
            document.getElementById(`sheet${sheetnumber}`).innerHTML += `<img src ='https://img.icons8.com/ultraviolet/32/000000/musical.png' style='width:50px;height:30px;position:relative;margin-left:3%;top:${(52-keynumber2)}%;'>`
        }else if(timelapse < beat4 && timelapse > beat8){
            document.getElementById(`sheet${sheetnumber}`).innerHTML += `<img src ='https://img.icons8.com/ios/32/000000/musical-filled.png' style='width:50px;height:30px;position:relative;margin-left:3%;top:${(52-keynumber2)}%;'>`
        }else if(timelapse < beat8){
            document.getElementById(`sheet${sheetnumber}`).innerHTML += `<img src ='https://img.icons8.com/material-sharp/32/000000/musical-notes.png' style='width:50px;height:30px;position:relative;margin-left:3%;top:${(52-keynumber2)}%;'>`
        };
*/
    };
};

