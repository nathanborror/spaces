package messages

import (
	"net/http"

	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/tokens"
	"github.com/nathanborror/spaces/dropbox"
	"github.com/nathanborror/spaces/rooms"
)

var repo = MessageSQLRepository("db.sqlite3")
var roomRepo = rooms.RoomSQLRepository("db.sqlite3")
var tokenRepo = tokens.TokenSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

func check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	au, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	authRepo.Ping(au) // Ping user

	hash := r.FormValue("hash")
	room := r.FormValue("room")
	text := r.FormValue("text")

	if hash == "" {
		hash = GenerateMessageHash(text)
	}

	m := &Message{Hash: hash, Room: room, User: au.Hash, Text: text}
	err = repo.Save(m)
	check(err, w)

	// Check for any resources in message
	dropbox.HandleDropboxFilesPut("DMX/Test.gdoc", text, r)

	// Push members
	go PushMembers(room, m.Text)

	// Redirect to message (this is kind of a hack so we return the right JSON
	// to the clients connected over websockets).
	http.Redirect(w, r, "/m/"+hash, http.StatusFound)
}
