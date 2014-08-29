package messages

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

// GenerateMessageHash returns a hash
func GenerateMessageHash(s string) (hash string) {
	time := time.Now().String()
	hasher := md5.New()
	io.WriteString(hasher, s+time)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
