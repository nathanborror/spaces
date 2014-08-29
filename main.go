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
	"github.com/nathanborror/gommon/hubspoke"
	"github.com/nathanborror/gommon/markdown"
	"github.com/nathanborror/gommon/render"
	"github.com/nathanborror/spaces/dropbox"
	"github.com/nathanborror/spaces/messages"
	"github.com/nathanborror/spaces/rooms"
)

var cookieStore = sessions.NewCookieStore([]byte("something-very-very-secret"))
var roomRepo = rooms.RoomSQLRepository("db.sqlite3")
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

// Room Handlers

func roomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	room, err := roomRepo.Load(hash)
	check(err, w)

	ml, err := messageRepo.List(hash)
	check(err, w)

	render.Render(w, r, "room", map[string]interface{}{
		"request":  r,
		"messages": ml,
		"room":     room,
	})
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	rooms, err := roomRepo.List(20)
	check(err, w)

	render.Render(w, r, "room_list", map[string]interface{}{
		"request": r,
		"rooms":   rooms,
	})
}

func roomFormHandler(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, "room_form", map[string]interface{}{
		"request": r,
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

var r = mux.NewRouter()

func main() {
	go hubspoke.Hub.Run()

	// Users
	r.HandleFunc("/login", auth.LoginHandler)
	r.HandleFunc("/logout", auth.LogoutHandler)
	r.HandleFunc("/register", auth.RegisterHandler)

	// Room
	r.HandleFunc("/r/create", auth.LoginRequired(roomFormHandler))
	r.HandleFunc("/r/save", auth.LoginRequired(rooms.SaveHandler))
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}", roomHandler)
	r.HandleFunc("/r/{hash:[a-zA-Z0-9-]+}/folder", roomFolderHandler)

	// Message
	r.HandleFunc("/m/save", auth.LoginRequired(messages.SaveHandler))
	r.HandleFunc("/m/{hash:[a-zA-Z0-9-]+}", messageHandler)

	// Dropbox
	http.HandleFunc("/dropbox", dropbox.HandleDropboxAuth)
	http.HandleFunc("/callback", dropbox.HandleDropboxCallback)

	r.HandleFunc("/ws", hubspoke.SpokeHandler)
	r.HandleFunc("/", auth.LoginRequired(roomsHandler))

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
