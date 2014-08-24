package messages

import (
	"net/http"
	"github.com/nathanborror/gommon/auth"
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
	parent := r.FormValue("parent")
	text := r.FormValue("text")

	if hash == "" {
		hash = GenerateMessageHash(text)
	}

	m := &Message{Hash: hash, Parent: parent, User: user.Hash, Text: text}
	err = repo.Save(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/m/"+hash, http.StatusFound)
}
