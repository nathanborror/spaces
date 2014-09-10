package rooms

import (
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"time"
)

// GenerateRoomHash returns a hash
func GenerateRoomHash(s string) (hash string) {
	time := time.Now().String()
	hasher := md5.New()
	io.WriteString(hasher, s+time)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// GenerateRoomMemberHash returns a hash
func GenerateRoomMemberHash(r string, u string) (hash string) {
	hasher := md5.New()
	io.WriteString(hasher, r+u)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// GenerateOneOnOneHash returns a hash for a one-on-one room
func GenerateOneOnOneHash(user1 string, user2 string) (hash string) {
	hasher := md5.New()

	// Sort users so hashes are in alphabetical order
	users := []string{user1, user2}
	sort.Strings(users)

	io.WriteString(hasher, users[0]+users[1])
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// JoinRoom joins a user with a room
func JoinRoom(room string, user string) error {
	hash := GenerateRoomMemberHash(room, user)
	rm := &RoomMember{Hash: hash, User: user, Room: room}
	err := roomMemberRepo.Save(rm)
	return err
}
