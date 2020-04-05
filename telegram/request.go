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

func KickChatMember(ctx context.Context, token string, chatID int64, userID int) error {
	// https://core.telegram.org/bots/api#kickchatmember
	const method = "kickChatMember"

	var (
		tgResp APIResponse
		url    string = buildURL(token, method)
	)

	var jsonStr = []byte(fmt.Sprintf(`{"chat_id":"%v","user_id":"%v"}`, chatID, userID))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
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
