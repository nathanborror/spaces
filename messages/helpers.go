package messages

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"
)

// GenerateMessageHash returns a hash
func GenerateMessageHash(s string) (hash string) {
	time := time.Now().String()
	hasher := md5.New()
	io.WriteString(hasher, s+time)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// FindStickers returns a set of MessageActiosn with the type 'sticker'
func FindStickers(text string) []*MessageAction {
	actions := []*MessageAction{}

	re := regexp.MustCompile(":(\\w+):")
	matches := re.FindAllStringSubmatch(text, -1)

	if len(matches) > 0 {
		for _, v := range matches {
			action := &MessageAction{Type: "sticker", Resource: v[1] + ".png", Raw: v[0]}
			actions = append(actions, action)
		}
	}

	return actions
}

// FindCommands returns a set of MessageActiosn with the type 'command'
func FindCommands(text string) []*MessageAction {
	actions := []*MessageAction{}

	re := regexp.MustCompile("/(join|msg|leave) (.+)?")
	match := re.FindStringSubmatch(text)

	if len(match) > 0 {
		action := &MessageAction{Type: match[1], Resource: match[2], Raw: match[0]}
		actions = append(actions, action)
	}

	return actions
}

// PushMembers sends APNS push notifications to the members of a room
func PushMembers(room string, text string) {
	members, err := roomRepo.ListMembers(room)
	if err != nil {
		log.Println("[Push]: ", err)
	}

	users := []string{}
	for _, m := range members {
		users = append(users, m.Hash)
	}

	tokenRepo.Push(users, text, "SpacesCert.pem", "SpacesKeyNoEnc.pem")
}
