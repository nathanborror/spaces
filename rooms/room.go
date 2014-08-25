package rooms

import "time"

// Room defines a blob
type Room struct {
	Hash     string    `json:"hash"`
	Users    string    `json:"users"`
	Name     string    `json:"name"`
	Folder   string    `json:"folder"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
