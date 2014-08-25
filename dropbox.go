package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// DateFormat is the format to use when decoding a time.
	DateFormat = time.RFC1123Z
)

// Token represents information about the oauth token
type Token struct {
	AccessToken string `json:"access_token"`
}

// Account represents information about the user account.
type Account struct {
	ReferralLink string `json:"referral_link,omitempty"` // URL for referral.
	DisplayName  string `json:"display_name,omitempty"`  // User name.
	UID          int    `json:"uid,omitempty"`           // User account ID.
	Country      string `json:"country,omitempty"`       // Country ISO code.
	QuotaInfo    struct {
		Shared int64 `json:"shared,omitempty"` // Quota for shared files.
		Quota  int64 `json:"quota,omitempty"`  // Quota in bytes.
		Normal int64 `json:"normal,omitempty"` // Quota for non-shared files.
	} `json:"quota_info"`
}

// SharedFolder represents information about a Dropbox shared folder
type SharedFolder struct {
	ID          string `json:"id,omitempty"`
	Path        string `json:"path,omitempty"`
	AccessLevel string `json:"access_level,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Members     struct {
		User struct {
			UID         int    `json:"uid,omitempty"`
			DisplayName string `json:"display_name,omitempty"`
			SameTeam    bool   `json:"same_team"`
		} `json:"user"`
		Role   string `json:"role,omitempty"`
		Active bool   `json:"active"`
	} `json:"members"`
}

// DBTime allow marshalling and unmarshalling of time.
type DBTime time.Time

// UnmarshalJSON unmarshals a time according to the Dropbox format.
func (dbt *DBTime) UnmarshalJSON(data []byte) error {
	var s string
	var err error
	var t time.Time

	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	if t, err = time.ParseInLocation(DateFormat, s, time.UTC); err != nil {
		return err
	}
	if t.IsZero() {
		*dbt = DBTime(time.Time{})
	} else {
		*dbt = DBTime(t)
	}
	return nil
}

// MarshalJSON marshals a time according to the Dropbox format.
func (dbt DBTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dbt).Format(DateFormat))
}

// Entry represents the metadata of a file or folder.
type Entry struct {
	Bytes       int     `json:"bytes,omitempty"`        // Size of the file in bytes.
	ClientMtime DBTime  `json:"client_mtime,omitempty"` // Modification time set by the client when added.
	Contents    []Entry `json:"contents,omitempty"`     // List of children for a directory.
	Hash        string  `json:"hash,omitempty"`         // Hash of this entry.
	Icon        string  `json:"icon,omitempty"`         // Name of the icon displayed for this entry.
	IsDeleted   bool    `json:"is_deleted,omitempty"`   // true if this entry was deleted.
	IsDir       bool    `json:"is_dir,omitempty"`       // true if this entry is a directory.
	MimeType    string  `json:"mime_type,omitempty"`    // MimeType of this entry.
	Modified    DBTime  `json:"modified,omitempty"`     // Date of last modification.
	Path        string  `json:"path,omitempty"`         // Absolute path of this entry.
	Revision    string  `json:"rev,omitempty"`          // Unique ID for this file revision.
	Root        string  `json:"root,omitempty"`         // dropbox or sandbox.
	Size        string  `json:"size,omitempty"`         // Size of the file humanized/localized.
	ThumbExists bool    `json:"thumb_exists,omitempty"` // true if a thumbnail is available for this entry.
}

func getCallbackURL(r *http.Request) string {
	scheme := "http"
	forwarded := r.Header.Get("X-Forwarded-Proto")
	if len(forwarded) > 0 {
		scheme = forwarded
	}
	return (&url.URL{
		Scheme: scheme,
		Host:   r.Host,
		Path:   "/callback",
	}).String()
}

// DecodeResponse returns a prepared interface to work with
func DecodeResponse(r *http.Response, m interface{}) {
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(m)
}

// Request allows you to make simple HTTP requests
func Request(url string, token string) (resp *http.Response, err error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	return resp, err
}

func handleDropboxAuth(w http.ResponseWriter, r *http.Request) {
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
			"client_id":     {appKey},
			"redirect_uri":  {getCallbackURL(r)},
			"response_type": {"code"},
			"state":         {csrf},
		}.Encode(), 302)
}

func handleDropboxCallback(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "csrf", MaxAge: -1})
	state := r.FormValue("state")
	cookie, _ := r.Cookie("csrf")
	if cookie == nil || cookie.Value != state {
		http.Error(w, "Possible CSRF attack.", http.StatusUnauthorized)
		return
	}

	resp, err := http.PostForm(fmt.Sprintf("https://%s:%s@api.dropbox.com/1/oauth2/token", appKey, appSecret),
		url.Values{
			"redirect_uri": {getCallbackURL(r)},
			"code":         {r.FormValue("code")},
			"grant_type":   {"authorization_code"},
		})
	if err != nil {
		panic(err)
	}

	var token Token
	DecodeResponse(resp, &token)

	// Saving session token
	session, _ := cookieStore.Get(r, "dropbox")
	session.Values["token"] = token.AccessToken
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}
