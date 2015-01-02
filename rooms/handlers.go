package rooms

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/crypto"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/spaces/dropbox"
)

var repo = RoomSQLRepository("db.sqlite3")
var roomMemberRepo = RoomMemberSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")
var cookieStore = sessions.NewCookieStore([]byte("something-very-very-secret"))

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := auth.GetAuthenticatedUser(r)

	hash := r.FormValue("hash")
	name := r.FormValue("name")
	kind := Open
	folder := r.FormValue("folder")
	created := time.Now()

	if hash == "" {
		hash = crypto.UniqueHash(name)
	}

	room, err := repo.Load(hash)
	if err == nil {
		created = room.Created
	}

	room = &Room{Hash: hash, Name: name, Kind: kind, Folder: folder, Created: created}
	err = repo.Save(room)
	render.Check(err, w)

	// Add members to room
	members := r.Form["members"]
	members = append(members, u.Key)
	for _, user := range members {
		hash = crypto.Hash(room.Hash + user)
		rm := &RoomMember{Hash: hash, User: user, Room: room.Hash}
		err = roomMemberRepo.Save(rm)
		render.Check(err, w)
	}

	http.Redirect(w, r, "/r/"+room.Hash, http.StatusFound)
}

// EditHandler handles editing of rooms
func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	u, _ := auth.GetAuthenticatedUser(r)

	users, err := userRepo.List(100)
	render.Check(err, w)

	room, err := repo.Load(hash)
	render.Check(err, w)

	members, err := roomMemberRepo.ListMembers(room.Hash)
	render.Check(err, w)

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
	render.Check(err, w)

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

	membership, err := roomMemberRepo.Load(hash, user.Key)
	render.Check(err, w)

	err = roomMemberRepo.Delete(membership.Hash)
	render.Check(err, w)

	http.Redirect(w, r, "/", 302)
}

// JoinHandler allows people to join rooms
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	user, _ := auth.GetAuthenticatedUser(r)

	room, err := repo.Load(hash)
	render.Check(err, w)

	err = JoinRoom(room.Hash, user.Key)
	render.Check(err, w)

	http.Redirect(w, r, "/r/"+room.Hash, 302)
}

// MemberHandler returns memebers for a room
func MemberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	room, err := repo.Load(hash)
	render.Check(err, w)

	members, err := roomMemberRepo.ListMembers(room.Hash)
	render.Check(err, w)

	render.Render(w, r, "room_members", map[string]interface{}{
		"request": r,
		"room":    room,
		"members": members,
	})
}

// ListHandler returns all available rooms
func ListHandler(w http.ResponseWriter, r *http.Request) {
	au, _ := auth.GetAuthenticatedUserKey(r)

	rooms, err := roomMemberRepo.ListRoomsForUser(au, 20)
	render.Check(err, w)

	joinable, err := roomMemberRepo.ListJoinableRoomsForUser(au, 20)
	render.Check(err, w)

	render.Render(w, r, "home", map[string]interface{}{
		"request":  r,
		"rooms":    rooms,
		"joinable": joinable,
	})
}

// FolderHandler returns a list of resources pertaining to the room
func FolderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	room, err := repo.Load(hash)
	render.Check(err, w)

	var folder dropbox.Entry

	if room.Folder != "" {
		session, _ := cookieStore.Get(r, "dropbox")
		token := fmt.Sprintf("%v", session.Values["token"])
		if token == "<nil>" { // FIXME: This could be more presentable
			http.Redirect(w, r, "/dropbox", 302)
		}

		url := fmt.Sprintf("https://api.dropbox.com/1/metadata/auto/%s", room.Folder)
		response, err := dropbox.Request("GET", url, token)
		render.Check(err, w)

		dropbox.DecodeResponse(response, &folder)
	}

	render.Render(w, r, "room_folder", map[string]interface{}{
		"request": r,
		"room":    room,
		"folder":  folder,
	})
}
