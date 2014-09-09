package rooms

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/render"
)

var repo = RoomSQLRepository("db.sqlite3")
var roomMemberRepo = RoomMemberSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

func check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.GetAuthenticatedUser(r)

	hash := r.FormValue("hash")
	name := r.FormValue("name")
	kind := Open
	folder := r.FormValue("folder")
	created := time.Now()

	if hash == "" {
		hash = GenerateRoomHash(name)
	}

	room, err := repo.Load(hash)
	if err == nil {
		created = room.Created
	}

	room = &Room{Hash: hash, Name: name, Kind: kind, Folder: folder, Created: created}
	err = repo.Save(room)
	check(err, w)

	// Add members to room
	members := r.Form["members"]
	members = append(members, u.Hash)
	for _, user := range members {
		hash = GenerateRoomMemberHash(room.Hash, user)
		rm := &RoomMember{Hash: hash, User: user, Room: room.Hash}
		err = roomMemberRepo.Save(rm)
		check(err, w)
	}

	http.Redirect(w, r, "/r/"+room.Hash, http.StatusFound)
}

// EditHandler handles editing of rooms
func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	u, _ := auth.GetAuthenticatedUser(r)

	users, err := userRepo.List(100)
	check(err, w)

	room, err := repo.Load(hash)
	check(err, w)

	members, err := repo.ListMembers(room.Hash)
	check(err, w)

	render.Render(w, r, "room_form", map[string]interface{}{
		"request": r,
		"room":    room,
		"users":   users,
		"members": members,
		"user":    u,
	})
}

// FormHandler presents a form for creating a new room
func FormHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.GetAuthenticatedUser(r)

	users, err := userRepo.List(100)
	check(err, w)

	render.Render(w, r, "room_form", map[string]interface{}{
		"request": r,
		"users":   users,
		"user":    u,
	})
}

// LeaveHandler allows people to leave rooms
func LeaveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	user, _ := auth.GetAuthenticatedUser(r)

	membership, err := roomMemberRepo.Load(hash, user.Hash)
	check(err, w)

	err = roomMemberRepo.Delete(membership.Hash)
	check(err, w)

	http.Redirect(w, r, "/", 302)
}

// JoinHandler allows people to join rooms
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	user, _ := auth.GetAuthenticatedUser(r)

	room, err := repo.Load(hash)
	check(err, w)

	hash = GenerateRoomMemberHash(room.Hash, user.Hash)
	rm := &RoomMember{Hash: hash, User: user.Hash, Room: room.Hash}
	err = roomMemberRepo.Save(rm)
	check(err, w)

	http.Redirect(w, r, "/r/"+room.Hash, 302)
}

// MemberHandler returns memebers for a room
func MemberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	room, err := repo.Load(hash)
	check(err, w)

	members, err := repo.ListMembers(room.Hash)
	check(err, w)

	render.Render(w, r, "room_members", map[string]interface{}{
		"request": r,
		"room":    room,
		"members": members,
	})
}

// ListHandler returns all available rooms
func ListHandler(w http.ResponseWriter, r *http.Request) {
	au, _ := auth.GetAuthenticatedUserHash(r)

	rooms, err := roomMemberRepo.List(au, 20)
	check(err, w)

	joinable, err := roomMemberRepo.ListJoinable(au, 20)
	check(err, w)

	render.Render(w, r, "room_list", map[string]interface{}{
		"request":  r,
		"rooms":    rooms,
		"joinable": joinable,
	})
}
