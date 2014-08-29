package rooms

import "time"

// Room defines a blob
type Room struct {
	Hash     string    `json:"hash"`
	Name     string    `json:"name"`
	Folder   string    `json:"folder"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// RoomMember defines a relationship between a User and a Room
type RoomMember struct {
	Hash    string    `json:"hash"`
	User    string    `json:"user"`
	Room    string    `json:"room"`
	Created time.Time `json:"created"`
}
