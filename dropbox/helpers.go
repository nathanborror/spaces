package dropbox

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

const (
	// DateFormat is the format to use when decoding a time.
	DateFormat = time.RFC1123Z
	// AppKey is the Dropbox App Key
	AppKey = "14l6emnb3m4jxye"
	// AppSecret is the Dropbox App Secret
	AppSecret = "8gdnanccsg7ty7f"
)

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
func Request(method string, url string, token string) (resp *http.Response, err error) {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	return resp, err
}
