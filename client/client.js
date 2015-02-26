"use strict";

var uuid = require("UUID");

function genteSock(str, debug){
	var ws = new WebSocket(str),
		ready = false,
		sendQueue = [],
		callbacks = {};

	ws.onopen = function(evt){
		if (debug) {
			console.log(evt);
		}

		ready = true;

		dequeueMessages();
	}

	ws.onmessage = function(evt){
		if (debug) {
			console.log(evt);
		}
	}

	function makeMsg(path, message){
		return {
			Id: uuid.V4(),
			Path: path,
			Message: message
		}
	}

	function dequeueMessages () {
		if (sendQueue.length === 0){
			return;
		}

		var msg = sendQueue.shift();

		ws.send(JSON.stringify(msg))

		if (sendQueue.length > 0){
			//Allow message to send, other threads to do work.
			window.setTimeout(dequeueMessages, 50)
		}
	}

	return {
		send: function(path, message, callback){
			var msg = makeMsg(path, message);
			sendQueue.push(msg)

			if (callback){
				callbacks[msg.Id] = callback;
			}

			if (ready){
				dequeueMessages();
			}
		},
	}
}


module.exports = {
	connect: function(str){
		return genteSock(str);
	}
}