package messages

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/anachronistic/apns"
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

// Push sends a push notification for a new (missed) message
func Push(text string, token string) {
	payload := apns.NewPayload()
	payload.Alert = text
	payload.Badge = 1
	payload.Sound = "bingbong.aiff"

	pn := apns.NewPushNotification()
	pn.DeviceToken = token
	pn.AddPayload(payload)

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "SpacesCert.pem", "SpacesKeyNoEnc.pem")
	resp := client.Send(pn)

	alert, _ := pn.PayloadString()
	if resp.Error != nil {
		log.Println("APNS Error: ", resp.Error)
	} else {
		log.Println("APNS Alert: ", alert)
		log.Println("APNS Success: ", resp.Success)
	}
}
