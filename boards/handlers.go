package boards

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/spaces/rooms"
)

var repo = BoardSQLRepository("db.sqlite3")
var roomRepo = rooms.RoomSQLRepository("db.sqlite3")
var pathRepo = PathSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

func check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SaveHandler saves a board
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.FormValue("hash")
	room := r.FormValue("room")
	name := r.FormValue("name")
	created := time.Now()

	if hash == "" {
		hash = GenerateHash()
	}

	board, err := repo.Load(hash)
	if err == nil {
		created = board.Created
	}

	board = &Board{Hash: hash, Room: room, Name: name, Created: created}
	err = repo.Save(board)
	check(err, w)

	http.Redirect(w, r, "/b/"+board.Hash, http.StatusFound)
}

// SavePathHandler saves a path for a board
func SavePathHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.GetAuthenticatedUser(r)

	hash := GenerateHash()
	board := r.FormValue("board")
	data := r.FormValue("data")

	path := &Path{Hash: hash, Board: board, Data: data, User: u.Hash}
	err := pathRepo.Save(path)
	check(err, w)

	http.Redirect(w, r, "/b/"+board, http.StatusFound)
}

// FormHandler presents a form for creating a board
func FormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["room"]

	room, err := roomRepo.Load(hash)
	check(err, w)

	render.Render(w, r, "board_form", map[string]interface{}{
		"request": r,
		"room":    room,
	})
}

// ListHandler returns all boards
func ListHandler(w http.ResponseWriter, r *http.Request) {
	boards, err := repo.List(20)
	check(err, w)

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
	check(err, w)

	paths, err := pathRepo.List(board.Hash)
	check(err, w)

	render.Render(w, r, "board", map[string]interface{}{
		"request": r,
		"board":   board,
		"paths":   paths,
	})
}

// DeletePathHandler removes a path from a board
func DeletePathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	err := pathRepo.Delete(hash)
	check(err, w)

	http.Redirect(w, r, "/", 302) // FIXME: Should redirect to the board
}

// DeleteAllPathHandler removes all paths from a board
func DeleteAllPathHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]

	err := pathRepo.DeleteAll(board)
	check(err, w)

	http.Redirect(w, r, "/", 302) // FIXME: Should redirect to the board
}
