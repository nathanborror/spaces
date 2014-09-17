package rooms

import "time"

// RoomMember defines a relationship between a User and a Room
type RoomMember struct {
	Hash    string    `json:"hash"`
	User    string    `json:"user"`
	Room    string    `json:"room"`
	Created time.Time `json:"created"`
}
