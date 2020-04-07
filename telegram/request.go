package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiURL = "https://api.telegram.org"

func buildURL(token string, method string) string {
	return fmt.Sprintf("%s/bot%s/%s", apiURL, token, method)
}

func doRequest(url string, reqBody []byte) error {
	var tgResp APIResponse

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &tgResp)
	if !tgResp.Ok {
		return fmt.Errorf("api error: %s, code: %v", tgResp.Description, tgResp.ErrorCode)
	}
	return nil
}

func SendMessage(ctx context.Context, token string, chatID int64, text string) error {
	// https://core.telegram.org/bots/api#sendmessage
	const method = "sendMessage"
	var (
		url  string = buildURL(token, method)
		body []byte = []byte(fmt.Sprintf(`{"chat_id":"%v","text":"%v","disable_web_page_preview":"true"}`, chatID, text))
	)
	return doRequest(url, body)
}

func Reply(ctx context.Context, token string, chatID int64, text string, messageID int) error {
	// https://core.telegram.org/bots/api#sendmessage
	const method = "sendMessage"
	var (
		url  string = buildURL(token, method)
		body []byte = []byte(fmt.Sprintf(`{"chat_id":"%v","text":"%v","reply_to_message_id":"%v","disable_web_page_preview":"true"}`, chatID, text, messageID))
	)
	return doRequest(url, body)
}

func KickChatMember(ctx context.Context, token string, chatID int64, userID int) error {
	// https://core.telegram.org/bots/api#kickchatmember
	const method = "kickChatMember"
	var (
		url  string = buildURL(token, method)
		body []byte = []byte(fmt.Sprintf(`{"chat_id":"%v","user_id":"%v"}`, chatID, userID))
	)
	return doRequest(url, body)
}

func UnbanChatMember(ctx context.Context, token string, chatID int64, userID int) error {
	// https://core.telegram.org/bots/api#unbanchatmember
	const method = "unbanChatMember"
	var (
		url  string = buildURL(token, method)
		body []byte = []byte(fmt.Sprintf(`{"chat_id":"%v","user_id":"%v"}`, chatID, userID))
	)
	return doRequest(url, body)
}
