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
	"github.com/nathanborror/gommon/markdown"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/gommon/spokes"
	"github.com/nathanborror/gommon/tokens"
	"github.com/nathanborror/spaces/dropbox"
	"github.com/nathanborror/spaces/messages"
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

func check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Rooms

func roomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	userHash, _ := auth.GetAuthenticatedUserHash(r)

	room, err := roomRepo.Load(hash)
	check(err, w)

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
	check(err, w)

	members, err := roomMemberRepo.ListMembers(room.Hash)
	check(err, w)

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
	user2, err := authRepo.Load(hash)
	check(err, w)

	// Create a one-on-one room between the logged in user and the clicked on
	// user if the room doesn't already exist.
	room, err := roomRepo.LoadOneOnOne(user1.Hash, user2.Hash)
	if err != nil {
		roomHash := rooms.GenerateOneOnOneHash(user1.Hash, user2.Hash)
		room = &rooms.Room{Hash: roomHash, Name: roomHash, Kind: rooms.OneOnOne, Folder: "", Created: time.Now()}
		err = roomRepo.Save(room)
		check(err, w)

		rooms.JoinRoom(roomHash, user1.Hash)
		rooms.JoinRoom(roomHash, user2.Hash)
	}

	render.Redirect(w, r, "/r/"+room.Hash)
}

var r = mux.NewRouter()

func main() {
	go spokes.Hub.Run()

	// Users
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/register", auth.RegisterHandler)
	r.HandleFunc("/u/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(oneOnOneHandler))

	// Tokens
	r.HandleFunc("/t/save", auth.LoginRequired(tokens.SaveHandler))

	// Room
	r.HandleFunc("/r/create", auth.LoginRequired(rooms.FormHandler))
	r.HandleFunc("/r/save", auth.LoginRequired(rooms.SaveHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(roomHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/edit", auth.LoginRequired(rooms.EditHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/folder", auth.LoginRequired(rooms.FolderHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/members", auth.LoginRequired(rooms.MemberHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/join", auth.LoginRequired(rooms.JoinHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/leave", auth.LoginRequired(rooms.LeaveHandler))
	r.HandleFunc("/r", auth.LoginRequired(rooms.ListHandler))

	// Message
	r.HandleFunc("/m/save", auth.LoginRequired(messages.SaveHandler))
	r.HandleFunc("/m/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(messages.MessageHandler))

	// Dropbox
	http.HandleFunc("/dropbox", dropbox.HandleDropboxAuth)
	http.HandleFunc("/callback", dropbox.HandleDropboxCallback)

	r.HandleFunc("/ws", spokes.SpokeHandler)
	r.HandleFunc("/", auth.LoginRequired(rooms.ListHandler))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
