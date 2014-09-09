package main

import (
	"fmt"
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

func memberOf(args ...interface{}) bool {
	if len(args) == 2 {
		user := args[1].(string)
		room := args[0].(string)

		_, err := roomMemberRepo.Load(room, user)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func init() {
	_ = render.RegisterTemplateFunction("markdown", markdown.Markdown)
	_ = render.RegisterTemplateFunction("ext", ext)
	_ = render.RegisterTemplateFunction("memberOf", memberOf)

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

// Room Handlers

func roomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	userHash, _ := auth.GetAuthenticatedUserHash(r)

	room, err := roomRepo.Load(hash)
	check(err, w)

	messages, err := messageRepo.List(hash, 20)
	check(err, w)

	isMember, _ := roomMemberRepo.Load(room.Hash, userHash)

	render.Render(w, r, "room", map[string]interface{}{
		"request":  r,
		"messages": messages,
		"room":     room,
		"isMember": isMember,
	})
}

func roomFolderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	room, err := roomRepo.Load(hash)
	check(err, w)

	var folder dropbox.Entry

	if room.Folder != "" {
		session, _ := cookieStore.Get(r, "dropbox")
		token := fmt.Sprintf("%v", session.Values["token"])
		if token == "<nil>" { // HACK
			http.Redirect(w, r, "/dropbox", 302)
		}

		url := fmt.Sprintf("https://api.dropbox.com/1/metadata/auto/%s", room.Folder)
		response, err := dropbox.Request("GET", url, token)
		if err != nil {
			panic(err)
		}

		dropbox.DecodeResponse(response, &folder)
	}

	render.Render(w, r, "room_folder", map[string]interface{}{
		"request": r,
		"room":    room,
		"folder":  folder,
	})
}

// Message Handlers

func messageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	message, err := messageRepo.Load(hash)
	check(err, w)

	room, err := roomRepo.Load(message.Room)
	check(err, w)

	user, err := authRepo.Load(message.User)
	check(err, w)

	render.Render(w, r, "message", map[string]interface{}{
		"request": r,
		"message": message,
		"room":    room,
		"user":    user,
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
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/folder", auth.LoginRequired(roomFolderHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/members", auth.LoginRequired(rooms.MemberHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/join", auth.LoginRequired(rooms.JoinHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/leave", auth.LoginRequired(rooms.LeaveHandler))
	r.HandleFunc("/r", auth.LoginRequired(rooms.ListHandler))

	// Message
	r.HandleFunc("/m/save", auth.LoginRequired(messages.SaveHandler))
	r.HandleFunc("/m/{hash:[a-zA-Z0-9-]+}", messageHandler)

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
