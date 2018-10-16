"use strict";
/// <reference path="jquery.d.ts" />
exports.__esModule = true;
var GUI_1 = require("./GUI");
$(function () {
    var socket = new GUI_1["default"]();
    socket.on('reply', function (msg) {
        $('#messages').append($('<li>').text(msg));
        return "callback from client " + msg;
    });
    $('form').submit(function () {
        socket.emit('msg', $('#m').val(), function (data) {
            $('#messages').append($('<li>').text('ACK CALLBACK: ' + data));
        });
        $('#m').val('');
        return false;
    });
    socket.connect();
});
