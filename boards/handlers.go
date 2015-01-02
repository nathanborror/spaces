package boards

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/crypto"
	"github.com/nathanborror/gommon/render"
)

var repo = BoardSQLRepository("db.sqlite3")
var pathRepo = PathSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

// SaveHandler saves a board
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("hash")
	room := r.FormValue("room")
	name := r.FormValue("name")
	created := time.Now()

	if hash == "" {
		hash = crypto.UniqueHash(name)
	}

	board, err := repo.Load(hash)
	if err == nil {
		created = board.Created
	}

	board = &Board{Hash: hash, Room: room, Name: name, Created: created}
	err = repo.Save(board)
	render.Check(err, w)

	http.Redirect(w, r, "/b/"+board.Hash, http.StatusFound)
}

// SavePathHandler saves a path for a board
func SavePathHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.GetAuthenticatedUser(r)

	hash := crypto.UniqueHash("")
	board := r.FormValue("board")
	data := r.FormValue("data")

	path := &Path{Hash: hash, Board: board, Data: data, User: u.Key}
	err := pathRepo.Save(path)
	render.Check(err, w)

	http.Redirect(w, r, "/b/"+board, http.StatusFound)
}

// FormHandler presents a form for creating a board
func FormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room := vars["room"]

	render.Render(w, r, "board_form", map[string]interface{}{
		"request": r,
		"room":    room,
	})
}

// ListHandler returns all boards
func ListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	boards, err := repo.List(hash)
	render.Check(err, w)

	render.Render(w, r, "board_list", map[string]interface{}{
		"request": r,
		"boards":  boards,
	})
}

// BoardHandler returns board with its paths
func BoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	board, err := repo.Load(hash)
	render.Check(err, w)

	paths, err := pathRepo.List(board.Hash)
	render.Check(err, w)

	render.Render(w, r, "board", map[string]interface{}{
		"request": r,
		"board":   board,
		"paths":   paths,
	})
}

// UndoPathHandler removes a path from a board
func UndoPathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	err := pathRepo.Delete(hash)
	render.Check(err, w)

	http.Redirect(w, r, "/", 302) // FIXME: Should redirect to the board
}

// ClearHandler removes all paths from a board
func ClearHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	err := repo.Clear(hash)
	render.Check(err, w)

	http.Redirect(w, r, "/b/"+hash, 302) // FIXME: Should redirect to the board
}
