
var handleMessage = function(data) {
  MessageManager.update(data);
};

var MessageManager = {};

// UpdateItems checks all the items on screen and adds any missing from the
// current dataset.
MessageManager.update = function(data) {
  var message_list = $('.ui-list');
  var current = _.map(message_list.find('> .ui-message'), function(obj) {
    return obj.id;
  });

  var hashes = _.map(data.messages, function(obj) {
    return obj.hash;
  });

  // Get the difference between the incoming hashes compared against
  // the existing hashes in the DOM.
  var diff = _.difference(hashes, current);

  // Insert any hashes that don't exist.
  for (var i=0; i<diff.length; i++) {
    var message = _.findWhere(data.messages, {'hash': diff[i]});
    var html = Message.html(message);
    Message.insert(html, message_list);
  }

  window.scrollTo(0, document.body.scrollHeight);
};

// Submit submits an message form.
MessageManager.submit = function(e) {
  e.preventDefault();
  var form = $(this);

  Message.save(form.serialize());

  form.find('input[name="text"]').val("");
};

var Message = {};

// HTML returns HTML necessary to render an message.
Message.html = function(message) {
  var text = message.text;

  if (message.text.slice(-1) == "?") {
    text = text.slice(0,-1);
  }

  switch (text) {
    case ":wolf:":
      text = '<img class="ui-message-sticker" src="/static/images/StickerWolf@2x.png">';
      break;
  }

  var html = $(''+
    '<div class="ui-message" id="'+message.hash+'">'+
      '<p><span class="ui-message-user">'+message.user+':</span> '+text+'</p>'+
    '</div>');

  html.data({
    'hash': message.hash,
    'parent': message.parent,
    'user': message.user,
    'text': message.text,
  });

  return html;
};

// Save saves a new message.
Message.save = function(data, complete) {
  $.post('/m/save', data, function(data) {
    if (complete) {
      complete(data);
    }
    window.SOCKET.request('/r/'+data.room.hash);
  }.bind(this));
};

// Insert adds message into a given room
Message.insert = function(message, room) {
  room.append(message);
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
});
