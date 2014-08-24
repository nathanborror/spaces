package rooms

import "time"

// Item defines a blob
type Room struct {
	Hash     string    `json:"hash"`
	Users    string    `json:"users"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
