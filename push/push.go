package push

import (
	"log"

	"github.com/nathanborror/gommon/tokens"
	"github.com/nathanborror/spaces/rooms"
)

var tokenRepo = tokens.TokenSQLRepository("db.sqlite3")
var roomMemberRepo = rooms.RoomMemberSQLRepository("db.sqlite3")

// PushMembers sends APNS push notifications to the members of a room
func PushMembers(room string, text string) {
	members, err := roomMemberRepo.ListMembers(room)
	if err != nil {
		log.Println("[Push]: ", err)
	}

	users := []string{}
	for _, m := range members {
		users = append(users, m.Hash)
	}

	tokenRepo.Push(users, text, "SpacesProdCert.pem", "SpacesProdKeyNoEnc.pem")
}
