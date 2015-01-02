package messages

import "regexp"

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
