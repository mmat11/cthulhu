package telegram

import (
	"fmt"
	"regexp"
	"strings"
)

const argsPattern = `("[^"]*"|[^"\s]+)(\s+|$)`

var argsPatternRe = regexp.MustCompile(argsPattern)

func GetUserName(user User) string {
	if user.UserName != "" {
		return fmt.Sprintf("@%s", user.UserName)
	}
	if user.LastName != "" {
		return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	return user.FirstName
}

func GetChatName(chat Chat) string {
	if chat.UserName != "" {
		return chat.UserName
	}
	return chat.Title
}

func CommandArgumentsSlice(args string) []string {
	lst := argsPatternRe.FindAllString(args, -1)
	for i, s := range lst {
		lst[i] = strings.TrimSpace(s)
	}
	return lst
}
