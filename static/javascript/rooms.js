/*
  depends on (
    underscore.js
    jquery.js
    users.js
  )
*/

var RoomManager = {};
var Room = {};

// Update checks all the items on screen and adds any missing from the
// current dataset.
RoomManager.update = function(data) {
  var room_list = $('.ui-list');
  var current = _.map(room_list.find('> .ui-row'), function(obj) {
    return obj.id;
  });

  var hashes = _.map(data.rooms, function(obj) {
    return obj.hash;
  });

  // Get the difference between the incoming hashes compared against
  // the existing hashes in the DOM.
  var diff = _.difference(hashes, current);

  // Insert any hashes that don't exist.
  for (var i=0; i<diff.length; i++) {
    var room = _.findWhere(data.rooms, {'hash': diff[i]});
    var html = Room.html(room);
    RoomManager.insert(html, room_list);
  }
};

// Insert adds rooms
RoomManager.insert = function(room, list) {
  list.append(room);
};

// HTML returns HTML necessary to render an room.
Room.html = function(room) {
  var recent = '';
  if (room.recent) {
    var user = User.get(room.recent.user);
    recent = user.name+': '+room.recent.text;
  }

  var html = $(''+
    '<div class="ui-row" id="'+room.hash+'">'+
      '<h4>'+
        '<a href="/r/'+room.hash+'">'+
          '<span class="ui-title">'+room.name+'</span>'+
          '<span class="ui-subtitle">'+recent+'</span>'+
        '</a>'+
      '</h4>'+
    '</div>');

  html.data({
    'hash': room.hash,
    'name': room.name,
  });

  return html;
};
