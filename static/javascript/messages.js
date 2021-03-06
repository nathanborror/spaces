/*
  depends on (
    underscore.js
    jquery.js
    users.js
  )
*/

var MessageManager = {};
var Message = {};

// UpdateItems checks all the items on screen and adds any missing from the
// current dataset.
MessageManager.update = function(data) {
  var message_list = $('.ui-detail-content');
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
    var user = User.get(message.user);

    var html = Message.html(message, user);
    MessageManager.insert(html, message_list);
  }

  window.scrollTo(0, document.body.scrollHeight);
};

// Submit submits an message form.
MessageManager.submit = function(e) {
  e.preventDefault();
  var form = $(this);

  Message.save(form.serialize());

  var textInput = form.find('input[name="text"]');

  // Clear inpupt field
  textInput.val("");
};

// Insert adds message into a given room
MessageManager.insert = function(message, room) {
  room.append(message);
};

// HTML returns HTML necessary to render an message.
Message.html = function(message, user) {
  var text = message.text;

  if (message.text.slice(-1) == "?") {
    text = text.slice(0,-1);
  }

  var re = /^(http|https):\/\/(.+).(png|jpg|gif)$/;
  if (text.search(re) != -1) {
    text = ''+
      '<div class="ui-message-images">'+
        '<a class="ui-message-image" href="'+text+'" target="_blank" style="background-image: url('+text+');"></a>'+
      '</div>';
  }

  // Actions
  var stickers = '';
  var command = '';

  for (var i=0; i<message.actions.length; i++) {
    var action = message.actions[i];
    switch(action.type) {
      case 'sticker': {
        stickers += '<span class="ui-message-sticker" style="background-image: url(/static/images/stickers/'+action.resource+');"></span>';
        break;
      }
      case 'join': {
        command = 'You joined '+action.resource;
        break;
      }
      case 'msg': {
        command = 'You messaged '+action.resource;
        break;
      }
      case 'leave': {
        command = 'You left '+action.resource;
        break;
      }
    }
  }

  if (stickers != '') {
    text = '<div class="ui-message-stickers">'+stickers+'</div>';
  }

  if (command != '') {
    var html = $(''+
      '<div class="ui-message ui-message-command" id="'+message.hash+'">'+
        '<p>'+command+'</p>'+
      '</div>');
  } else {
    text = text.replace(/@(\w+)/, '<a href="/u/$1" class="ui-mention"><span>@</span>$1</a>');
    var html = $(''+
      '<div class="ui-message" id="'+message.hash+'">'+
        '<a class="ui-message-profile" href="#"></a>'+
        '<p><a href="/u/'+user.key+'" class="ui-message-user">'+user.name+'</a> '+text+'</p>'+
      '</div>');
  }

  html.data({
    'hash': message.hash,
    'room': message.room,
    'user': message.user,
    'text': message.text,
  });

  return html;
};

// Save saves a new message.
Message.save = function(data, complete) {
  $.ajax({
    type: 'POST',
    url: '/m/save',
    data: data,
    success: function(data) {
      if (complete) {
        complete(data);
      }
      window.SOCKET.request('/r/'+data.message.room);
    }.bind(this),
    error: function(xhr, status, error) {
      alert('There was '+status+' when trying to send this message. Please contact nathan@dropbox.com.');
    }
  });
};
