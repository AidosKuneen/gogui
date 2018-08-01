/// <reference path="jquery.d.ts" />

import GUI from "./GUI";

$(function () {
    const socket = new GUI();
    socket.on('reply', (msg:string) =>{
      $('#messages').append($('<li>').text(msg));
      return "callback from client "+msg;
    });
    $('form').submit(function () {
      socket.emit('msg', $('#m').val(), (data:string)=> {
        $('#messages').append($('<li>').text('ACK CALLBACK: ' + data));
      });
      $('#m').val('');
      return false;
    });
    socket.connect()
  })