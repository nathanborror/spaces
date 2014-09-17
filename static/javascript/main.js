/*
  depends on (
    socket.js
    underscore.js
    jquery.js
    rooms.js
    messages.js
  )
*/

var handleMessage = function(data) {
  MessageManager.update(data);
  RoomManager.update(data);
  console.log("[WebSocket]: Received message");
};

// HACK
$(function() {
  var body = $('body');

  // Rooms
  body.on('submit', '.ui-room-form', MessageManager.submit);

  // Input fields
  $('.ux-focus').focus();

  // WebSocket
  window.SOCKET.onclose = function(e) {
    $("#alert").append("<p>Oops! You've been disconnected. <a href='javascript:location.reload();'>Reload</a> to fix this.</p>");
  }

  window.SOCKET.onopen = function(e) {
    $("#alert >").remove();
    window.SOCKET.subscribe(window.location.pathname, handleMessage);
    window.SOCKET.request(window.location.pathname);
  }

  window.scrollTo(0, document.body.scrollHeight);

  UserManager.init();
});
