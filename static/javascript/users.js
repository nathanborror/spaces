/*
  depends on (
    underscore.js
    jquery.js
  )
*/

var User = {};

User.list = [];

User.get = function(hash) {
  return _.findWhere(User.list, {'hash': hash})
}

var UserManager = {};

UserManager.init = function() {
  $.ajax({
    type: 'GET',
    url: '/u',
    success: function(data) {
      User.list = data.users;
    }.bind(this),
    error: function(xhr, status, error) {
      alert('There was '+status+' when trying to retrieve list of users.');
    }
  });
}
