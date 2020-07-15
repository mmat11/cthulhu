package telegram

import "fmt"

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
