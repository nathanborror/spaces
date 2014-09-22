
var kActionRequest = 'request';
var kActionSubscribe = 'subscribe';
var kPort = '8082';

window.SOCKET = new WebSocket("ws://"+window.location.host+"/ws");

window.SOCKET.subscriptions = {};

window.SOCKET.onopen = function(e) {
  console.log('[WebSocket]: Connection opened.');
};

window.SOCKET.onclose = function(e) {
  console.log('[WebSocket]: Connection closed.');
};

window.SOCKET.onmessage = function(e) {
  if (!e.data) return;

  var message = JSON.parse(e.data);
  var channel = message.Channel;
  var data = message.Data;

  callbacks = window.SOCKET.subscriptions[channel] || [];
  callbacks.forEach(function(cb) {
    cb.call(e, data);
  });
};

// on registers a callback when a new message on channel `channel` occurs.
window.SOCKET.subscribe = function(channel, callback) {
  var payload = JSON.stringify({'url': channel, 'action': kActionSubscribe, 'port': kPort});
  if (!window.SOCKET.subscriptions[channel]) {
    window.SOCKET.subscriptions[channel] = [];
  }
  window.SOCKET.subscriptions[channel].push(callback);
  window.SOCKET.send(payload);
};

window.SOCKET.unsubscribe = function(channel) {
  window.SOCKET.subscriptions[channel].pop();
}

window.SOCKET.unsubscribeAll = function() {
  window.SOCKET.subscriptions = {};
}

window.SOCKET.request = function(url) {
  var payload = JSON.stringify({'url': url, 'action': kActionRequest, 'port': kPort});
  window.SOCKET.send(payload);
};
