package devices

import "time"

// Device defines a message sent from a user
type Device struct {
	Token    string    `json:"token"`
	Make     string    `json:"make"`
	User     string    `json:"user"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}
