package rooms

import (
	"sort"

	"github.com/nathanborror/gommon/crypto"
)

// GenerateOneOnOneHash returns a hash for a one-on-one room
func GenerateOneOnOneHash(user1 string, user2 string) (hash string) {
	// Sort users so hashes are in alphabetical order
	users := []string{user1, user2}
	sort.Strings(users)

	return crypto.Hash(users[0] + users[1])
}

// JoinRoom joins a user with a room
func JoinRoom(room string, user string) error {
	hash := crypto.Hash(room + user)
	rm := &RoomMember{Hash: hash, User: user, Room: room}
	err := roomMemberRepo.Save(rm)
	return err
}
