package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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

	messages, err := messageRepo.List(hash, 20)
	check(err, w)

	members, err := roomMemberRepo.ListMembers(room.Hash)
	check(err, w)

	isMember, _ := roomMemberRepo.Load(room.Hash, userHash)

	render.Render(w, r, "room", map[string]interface{}{
		"request":  r,
		"messages": messages,
		"room":     room,
		"members":  members,
		"isMember": isMember,
	})
}

// Users

func usersHandler(w http.ResponseWriter, r *http.Request) {
	au, _ := auth.GetAuthenticatedUser(r)

	users, err := authRepo.List(100)
	check(err, w)

	render.Render(w, r, "user_list", map[string]interface{}{
		"request":  r,
		"authUser": au,
		"users":    users,
	})
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	au, _ := auth.GetAuthenticatedUser(r)

	user, err := authRepo.Load(hash)
	check(err, w)

	// TODO:
	// Create a OneOnOne thread with this person if it doesn't already exist

	render.Render(w, r, "user", map[string]interface{}{
		"request":  r,
		"authUser": au,
		"user":     user,
	})
}

var r = mux.NewRouter()

func main() {
	go spokes.Hub.Run()

	// Users
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/register", auth.RegisterHandler)
	r.HandleFunc("/u", auth.LoginRequired(usersHandler))
	r.HandleFunc("/u/{hash:[a-zA-Z0-9-]+}", auth.LoginRequired(userHandler))

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
