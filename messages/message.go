package messages

import "time"

// Message defines a message sent from a user
type Message struct {
	Hash     string    `json:"hash"`
	Room     string    `json:"room"`
	User     string    `json:"user"`
	Text     string    `json:"text"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
