package messages

import (
	"net/http"

	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/spaces/dropbox"
)

var repo = MessageSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	hash := r.FormValue("hash")
	room := r.FormValue("room")
	text := r.FormValue("text")

	if hash == "" {
		hash = GenerateMessageHash(text)
	}

	m := &Message{Hash: hash, Room: room, User: user.Hash, Text: text}
	err = repo.Save(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for any resources in message
	dropbox.HandleDropboxFilesPut("DMX/Test.gdoc", text, r)

	// Redirect to message (this is kind of a hack so we return the right JSON
	// to the clients connected over websockets).
	http.Redirect(w, r, "/m/"+hash, http.StatusFound)
}
