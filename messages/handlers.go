package messages

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/crypto"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/spaces/dropbox"
)

var repo = MessageSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	au, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	hash := r.FormValue("hash")
	room := r.FormValue("room")
	text := r.FormValue("text")

	if hash == "" {
		hash = crypto.UniqueHash(text)
	}

	m := &Message{Hash: hash, Room: room, User: au.Key, Text: text}
	err = repo.Save(m)
	render.Check(err, w)

	// Check for any resources in message
	go dropbox.HandleDropboxFilesPut("DMX/Test.gdoc", text, r)

	// Redirect to message (this is kind of a hack so we return the right JSON
	// to the clients connected over websockets).
	http.Redirect(w, r, "/m/"+hash, http.StatusFound)
}

// MessageHandler returns a message
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	message, err := repo.Load(hash)
	render.Check(err, w)

	user, err := authRepo.Get(message.User)
	render.Check(err, w)

	render.Render(w, r, "message", map[string]interface{}{
		"request": r,
		"message": message,
		"user":    user,
	})
}
