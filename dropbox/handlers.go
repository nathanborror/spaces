package dropbox

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/gorilla/sessions"
)

var cookieStore = sessions.NewCookieStore([]byte("something-very-very-secret"))

func check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleDropboxFilesPut handles a Dropbox file put request
func HandleDropboxFilesPut(name string, text string, r *http.Request) {
	// Check to see if a valid resource exists
	re, _ := regexp.Compile("(http://|https://)docs.google.com/a/dropbox.com/(document|spreadsheets|presentation)/d/(.+)")
	match := re.FindStringSubmatch(text)

	if len(match) == 0 {
		log.Println("Dropbox: No resource found in message")
		return
	}

	session, _ := cookieStore.Get(r, "dropbox")
	token := fmt.Sprintf("%v", session.Values["token"])
	if token == "<nil>" { // HACK
		// TODO: do something here
	}

	content := fmt.Sprintf("{\"url\": \"%s\", \"resource_id\": \"%s\"}", match[0], match[3])
	url := fmt.Sprintf("https://api-content.dropbox.com/1/files_put/auto/%s", name)
	size := int64(len(content))

	req, _ := http.NewRequest("PUT", url, bytes.NewBufferString(content))
	req.Header.Set("Content-Length", strconv.FormatInt(size, 10))
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	var entry Entry
	DecodeResponse(response, &entry)
}

// HandleDropboxAuth handles Dropbox authentication using OAuth2
func HandleDropboxAuth(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dropbox" {
		http.NotFound(w, r)
		return
	}

	b := make([]byte, 18)
	rand.Read(b)
	csrf := base64.StdEncoding.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{Name: "csrf", Value: csrf})

	http.Redirect(w, r, "https://www.dropbox.com/1/oauth2/authorize?"+
		url.Values{
			"client_id":     {AppKey},
			"redirect_uri":  {AppCallback},
			"response_type": {"code"},
			"state":         {csrf},
		}.Encode(), 302)
}

// HandleDropboxCallback handles Dropbox OAuth2 callback
func HandleDropboxCallback(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "csrf", MaxAge: -1})
	state := r.FormValue("state")
	cookie, _ := r.Cookie("csrf")
	if cookie == nil || cookie.Value != state {
		http.Error(w, "Possible CSRF attack.", http.StatusUnauthorized)
		return
	}

	resp, err := http.PostForm(fmt.Sprintf("https://%s:%s@api.dropbox.com/1/oauth2/token", AppKey, AppSecret),
		url.Values{
			"redirect_uri": {AppCallback},
			"code":         {r.FormValue("code")},
			"grant_type":   {"authorization_code"},
		})
	check(err, w)

	var token Token
	DecodeResponse(resp, &token)

	// Saving session token
	session, _ := cookieStore.Get(r, "dropbox")
	session.Values["token"] = token.AccessToken
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
