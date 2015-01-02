/*
  depends on (
    underscore.js
    jquery.js
  )
*/

var User = {};

User.list = [];

User.get = function(key) {
  return _.findWhere(User.list, {'key': key})
}

var UserManager = {};

UserManager.init = function(callback) {
  $.ajax({
    type: 'GET',
    url: '/u',
    success: function(data) {
      User.list = data.users;
      callback();
    }.bind(this),
    error: function(xhr, status, error) {
      alert('There was '+status+' when trying to retrieve list of users.');
    }
  });
}
