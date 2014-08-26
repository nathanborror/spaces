package messages

import (
	"time"

	"github.com/nathanborror/gommon/auth"
)

var authRepo = auth.AuthSQLRepository("db.sqlite3")

// MarshalPreparable can supply an alternative value in preparation for marshalling
type MarshalPreparable interface {
	MarshalPrepare() interface{}
}

// Message defines a message sent from a user
type Message struct {
	Hash     string    `json:"hash"`
	Room     string    `json:"room"`
	User     string    `json:"user"`
	Text     string    `json:"text"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// MessageList returns a list of Messages
type MessageList []*Message

// UserObject returns the actual User instance
func (m Message) UserObject() *auth.User {
	u, err := authRepo.Load(m.User)
	if err != nil {
		return nil
	}
	return u
}

// MarshalPrepare output
func (m Message) MarshalPrepare() interface{} {
	return struct {
		Message
		User *auth.User `json:"user"`
	}{m, m.UserObject()}
}

// MarshalPrepare prepares a list of messages
func (ml MessageList) MarshalPrepare() interface{} {
	result := make([]interface{}, 0, len(ml))
	for _, m := range ml {
		result = append(result, m.MarshalPrepare())
	}
	return result
}
