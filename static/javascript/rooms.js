/*
  depends on (
    underscore.js
    jquery.js
    users.js
  )
*/

var RoomManager = {};
var Room = {};

RoomManager.init = function() {
  var body = $('body');
  body.on('click', '.ui-row-room a', Room.handleClick);

  body.on('click', '.ui-boards', function(e) {
    window.location = '/r/'+window.ACTIVE_ROOM.hash+'/boards';
  });

  body.on('click', '.ui-folder', function(e) {
    window.location = '/r/'+window.ACTIVE_ROOM.hash+'/folder';
  });

  body.on('click', '.ui-members', function(e) {
    window.location = '/r/'+window.ACTIVE_ROOM.hash+'/members';
  });
};

// Update checks all the items on screen and adds any missing from the
// current dataset.
RoomManager.update = function(data) {
  var room_list = $('.ui-master-content');
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

// Make room active
RoomManager.makeActive = function(hash, name) {
  window.ACTIVE_ROOM = {'hash': hash, 'name': name};

  // Change url
  window.history.pushState({}, name, '/r/'+hash)

  // Make body class active
  $('body').addClass('room-active');

  // Highlight the room in the list
  $('.ui-row-selected').removeClass('ui-row-selected');
  $('#'+hash).addClass('ui-row-selected');

  // Change title
  $('header .ui-title').html(name);

  // Set room input
  $('.ui-room-form input[name="room"]').val(hash);

  // Remove any prior room messages
  $('.ui-detail-content').html('');

  // Subscribe to content
  window.SOCKET.subscribe('/r/'+hash, handleMessage);
  window.SOCKET.request('/r/'+hash);
}

// HTML returns HTML necessary to render an room.
Room.html = function(room) {
  var url = '/r/'+room.hash;

  var recent = '';
  if (room.recent) {
    var user = User.get(room.recent.user);
    recent = user.name+': '+room.recent.text;
  }

  var name = room.name;
  if (room.kind == 'oneonone') {
    switch (room.members.length) {
      case 1:
        name = "Note to self";
        break;
      case 2:
        var user1 = User.get(room.members[0].user);
        var user2 = User.get(room.members[1].user);
        name = user1.name+", "+user2.name;
        break;
    }
  }

  var html = $(''+
    '<div class="ui-row ui-row-room" id="'+room.hash+'">'+
      '<h4>'+
        '<a href="'+url+'">'+
          '<span class="ui-title">'+name+'</span>'+
          '<span class="ui-subtitle">'+recent+'</span>'+
        '</a>'+
      '</h4>'+
      '<a class="ui-boards" href="'+url+'/boards">Boards</a>'+
      '<a class="ui-folder" href="'+url+'/folder">Folder</a>'+
      '<a class="ui-members" href="'+url+'/members">Members</a>'+
    '</div>');

  html.data({
    'hash': room.hash,
    'name': room.name,
  });

  if (window.location.pathname == url) {
    RoomManager.makeActive(room.hash, room.name);
    html.addClass('ui-row-selected');
  }

  return html;
};

Room.handleClick = function(e) {
  e.preventDefault();

  var body = $('body');
  var room = $(this).parents('.ui-row');

  // Unsubscribe from previous room
  if (window.ACTIVE_ROOM) {
      window.SOCKET.unsubscribe('/r/'+window.ACTIVE_ROOM.hash);
  }

  RoomManager.makeActive(room.data('hash'), room.data('name'));
};
