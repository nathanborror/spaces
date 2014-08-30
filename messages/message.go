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

// MessageAction defines an action thats encoded in the message (e.g. stickers, room events, etc.)
type MessageAction struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
	Raw      string `json:"raw"`
}

// GetUser returns the actual User instance
func (m Message) GetUser() *auth.User {
	u, err := authRepo.Load(m.User)
	if err != nil {
		return nil
	}
	return u
}

// GetActions returns actions (e.g. stickers, joins, etc.)
func (m Message) GetActions() []*MessageAction {
	commands := FindCommands(m.Text)
	if len(commands) > 0 {
		return commands
	}

	stickers := FindStickers(m.Text)
	return stickers
}

// MarshalPrepare output
func (m Message) MarshalPrepare() interface{} {
	return struct {
		Message
		User    *auth.User       `json:"user"`
		Actions []*MessageAction `json:"actions"`
	}{m, m.GetUser(), m.GetActions()}
}

// MarshalPrepare prepares a list of messages
func (ml MessageList) MarshalPrepare() interface{} {
	result := make([]interface{}, 0, len(ml))
	for _, m := range ml {
		result = append(result, m.MarshalPrepare())
	}
	return result
}
