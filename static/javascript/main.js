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
  console.log("[WebSocket]: Received message");

  MessageManager.update(data);
  RoomManager.update(data);
};

window.ACTIVE_ROOM;

// HACK
$(function() {
  var body = $('body');

  // Rooms
  body.on('submit', '.ui-room-form', MessageManager.submit);

  // WebSocket
  window.SOCKET.onclose = function(e) {
    $("#alert").append("<p>Oops! You've been disconnected. <a href='javascript:location.reload();'>Reload</a> to fix this.</p>");
  }

  window.SOCKET.onopen = function(e) {
    $("#alert >").remove();

    UserManager.init(function() {
      // Always subscribe for the list of rooms
      window.SOCKET.subscribe('/', handleMessage);
      window.SOCKET.request('/');

      RoomManager.init();

      // Input fields
      $('.ux-focus').focus();
    });

  }
});
