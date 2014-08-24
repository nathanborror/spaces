package rooms

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

// GenerateItemHash returns a hash
func GenerateItemHash(s string) (hash string) {
	time := time.Now().String()
	hasher := md5.New()
	io.WriteString(hasher, s+time)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
