package boards

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

// GenerateHash returns a hash
func GenerateHash() (hash string) {
	time := time.Now().String()
	hasher := md5.New()
	io.WriteString(hasher, time)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
