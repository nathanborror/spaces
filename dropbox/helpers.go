package dropbox

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	// DateFormat is the format to use when decoding a time.
	DateFormat = time.RFC1123Z
	// AppKey is the Dropbox App Key
	AppKey = "14l6emnb3m4jxye"
	// AppSecret is the Dropbox App Secret
	AppSecret = "8gdnanccsg7ty7f"
	// AppCallback is the callback url
	AppCallback = "https://spaces.nthn.me/callback"
)

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
