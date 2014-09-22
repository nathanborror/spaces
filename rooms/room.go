package rooms

import (
	"time"

	"github.com/nathanborror/spaces/messages"
)

var messageRepo = messages.MessageSQLRepository("db.sqlite3")

// OneOnOne is a room with two people only. Open is a room anyone can join.
// Closed is only joinable upon request.
const (
	OneOnOne string = "oneonone"
	Open     string = "open"
	Closed   string = "closed"
	Secret   string = "secret"
)

// MarshalPreparable can supply an alternative value in preparation for marshalling
type MarshalPreparable interface {
	MarshalPrepare() interface{}
}

// Room defines a blob
type Room struct {
	Hash     string    `json:"hash"`
	Name     string    `json:"name"`
	Kind     string    `json:"kind"`
	Folder   string    `json:"folder"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// RoomList returns a list of rooms with marshalable annotations
type RoomList []*Room

// GetRecents returns actions (e.g. stickers, joins, etc.)
func (r Room) GetRecent() *messages.Message {
	ml, err := messageRepo.List(r.Hash, 1)
	if err != nil {
		return nil
	}
	if len(ml) < 1 {
		return nil
	}
	return ml[0]
}

func (r Room) GetMembers() []*RoomMember {
	m, err := roomMemberRepo.List(r.Hash)
	if err != nil {
		return nil
	}
	return m
}

// MarshalPrepare output
func (r Room) MarshalPrepare() interface{} {
	return struct {
		Room
		Recent *messages.Message `json:"recent"`
		Members []*RoomMember `json:"members"`
	}{r, r.GetRecent(), r.GetMembers()}
}

// MarshalPrepare prepares a list of rooms
func (rl RoomList) MarshalPrepare() interface{} {
	result := make([]interface{}, 0, len(rl))
	for _, r := range rl {
		result = append(result, r.MarshalPrepare())
	}
	return result
}
