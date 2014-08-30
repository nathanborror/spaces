package rooms

import "time"

// OneOnOne is a room with two people only. Open is a room anyone can join.
// Closed is only joinable upon request.
const (
	OneOnOne string = "oneonone"
	Open     string = "open"
	Closed   string = "closed"
	Secret   string = "secret"
)

// Room defines a blob
type Room struct {
	Hash     string    `json:"hash"`
	Name     string    `json:"name"`
	Kind     string    `json:"kind"`
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
