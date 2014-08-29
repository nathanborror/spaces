package rooms

import (
	"crypto/md5"
	"fmt"
	"io"
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
