package rooms

import (
	"net/http"
	"time"

	"github.com/nathanborror/gommon/auth"
)

var repo = RoomSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("hash")
	name := r.FormValue("name")
	users := r.FormValue("users")
	folder := r.FormValue("folder")
	created := time.Now()

	if hash == "" {
		hash = GenerateItemHash(name)
	}

	room, err := repo.Load(hash)
	if err == nil {
		created = room.Created
	}

	room = &Room{Hash: hash, Users: users, Name: name, Folder: folder, Created: created}
	err = repo.Save(room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/r/"+hash, http.StatusFound)
}
