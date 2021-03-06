package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathanborror/gommon/auth"
	"github.com/nathanborror/gommon/crypto"
	"github.com/nathanborror/gommon/markdown"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/gommon/spokes"
	"github.com/nathanborror/gommon/tokens"
	"github.com/nathanborror/spaces/boards"
	"github.com/nathanborror/spaces/dropbox"
	"github.com/nathanborror/spaces/messages"
	"github.com/nathanborror/spaces/push"
	"github.com/nathanborror/spaces/rooms"
)

var cookieStore = sessions.NewCookieStore([]byte("something-very-very-secret"))
var roomRepo = rooms.RoomSQLRepository("db.sqlite3")
var roomMemberRepo = rooms.RoomMemberSQLRepository("db.sqlite3")
var messageRepo = messages.MessageSQLRepository("db.sqlite3")
var authRepo = auth.AuthSQLRepository("db.sqlite3")

func ext(name string) string {
	e := strings.Split(name, ".")
	if len(e) > 1 {
		return e[len(e)-1]
	}
	return ""
}

func init() {
	_ = render.RegisterTemplateFunction("markdown", markdown.Markdown)
	_ = render.RegisterTemplateFunction("ext", ext)

	cookieStore.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

// Rooms

func roomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	userHash, _ := auth.GetAuthenticatedUserKey(r)

	room, err := roomRepo.Load(hash)
	render.Check(err, w)

	// If this is a private one-on-one room check whether the logged
	// in user is part of that room. If not, then they shouldn't be
	// able to view this room.
	if room.Kind == rooms.OneOnOne {
		_, err := roomMemberRepo.Load(room.Hash, userHash)
		if err != nil {
			render.Redirect(w, r, "/")
			return
		}
	}

	messages, err := messageRepo.List(hash, 20)
	render.Check(err, w)

	members, err := roomMemberRepo.ListMembers(room.Hash)
	render.Check(err, w)

	isMember, _ := roomMemberRepo.Load(room.Hash, userHash)

	render.Render(w, r, "room", map[string]interface{}{
		"request":  r,
		"room":     room,
		"messages": messages,
		"members":  members,
		"isMember": isMember,
	})
}

// Users

func oneOnOneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	user1, _ := auth.GetAuthenticatedUser(r)
	user2, err := authRepo.Get(hash)
	render.Check(err, w)

	// Create a one-on-one room between the logged in user and the clicked on
	// user if the room doesn't already exist.
	room, err := roomRepo.LoadOneOnOne(user1.Key, user2.Key)
	if err != nil {
		roomHash := rooms.GenerateOneOnOneHash(user1.Key, user2.Key)
		room = &rooms.Room{Hash: roomHash, Name: roomHash, Kind: rooms.OneOnOne, Folder: "", Created: time.Now()}
		err = roomRepo.Save(room)
		render.Check(err, w)

		rooms.JoinRoom(roomHash, user1.Key)
		rooms.JoinRoom(roomHash, user2.Key)
	}

	render.Redirect(w, r, "/r/"+room.Hash)
}

func userListHandler(w http.ResponseWriter, r *http.Request) {
	users, err := authRepo.List(100)
	render.Check(err, w)

	render.Render(w, r, "user_list", map[string]interface{}{
		"request": r,
		"users":   users,
	})
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	au, _ := auth.GetAuthenticatedUserKey(r)

	rs, err := roomMemberRepo.ListRoomsForUser(au, 20)
	render.Check(err, w)

	joinable, err := roomMemberRepo.ListJoinableRoomsForUser(au, 20)
	render.Check(err, w)

	render.Render(w, r, "home", map[string]interface{}{
		"request":  r,
		"rooms":    rs,
		"joinable": joinable,
	})
}

func messageSaveHandler(w http.ResponseWriter, r *http.Request) {
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

	m := &messages.Message{Hash: hash, Room: room, User: au.Key, Text: text}
	err = messageRepo.Save(m)
	render.Check(err, w)

	// Check for any resources in message
	go dropbox.HandleDropboxFilesPut("DMX/Test.gdoc", text, r)

	// Push members
	go push.PushMembers(room, m.Text)

	// Redirect to message (this is kind of a hack so we return the right JSON
	// to the clients connected over websockets).
	http.Redirect(w, r, "/m/"+hash, http.StatusFound)
}

var r = mux.NewRouter()

func main() {
	go spokes.Hub.Run()

	// Users
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/register", auth.RegisterHandler)
	r.HandleFunc("/u/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(oneOnOneHandler))
	r.HandleFunc("/u", auth.LoginRequired(userListHandler))

	// Tokens
	r.HandleFunc("/t/save", auth.LoginRequired(tokens.SaveHandler))

	// Boards
	r.HandleFunc("/b/save", auth.LoginRequired(boards.SaveHandler))
	r.HandleFunc("/b/{room:[a-zA-Z0-9-]+}/create", auth.LoginRequired(boards.FormHandler))
	r.HandleFunc("/b/{hash:[a-zA-Z0-9-]+}/clear", auth.LoginRequired(boards.ClearHandler))
	r.HandleFunc("/b/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(boards.BoardHandler))

	// Paths
	r.HandleFunc("/p/save", auth.LoginRequired(boards.SavePathHandler))
	r.HandleFunc("/p/{hash:[a-zA-Z0-9-]+}/undo", auth.LoginRequired(boards.UndoPathHandler))

	// Room
	r.HandleFunc("/r/create", auth.LoginRequired(rooms.FormHandler))
	r.HandleFunc("/r/save", auth.LoginRequired(rooms.SaveHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(roomHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/edit", auth.LoginRequired(rooms.EditHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/folder", auth.LoginRequired(rooms.FolderHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/members", auth.LoginRequired(rooms.MemberHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/boards", auth.LoginRequired(boards.ListHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/join", auth.LoginRequired(rooms.JoinHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/leave", auth.LoginRequired(rooms.LeaveHandler))
	r.HandleFunc("/r", auth.LoginRequired(roomsHandler))

	// Message
	r.HandleFunc("/m/save", auth.LoginRequired(messageSaveHandler))
	r.HandleFunc("/m/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(messages.MessageHandler))

	// Dropbox
	http.HandleFunc("/dropbox", dropbox.HandleDropboxAuth)
	http.HandleFunc("/callback", dropbox.HandleDropboxCallback)

	r.HandleFunc("/ws", spokes.SpokeHandler)
	r.HandleFunc("/", auth.LoginRequired(roomsHandler))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
