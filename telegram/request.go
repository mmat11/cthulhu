package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log"
)

type Service interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
	Reply(ctx context.Context, chatID int64, text string, messageID int) error
	KickChatMember(ctx context.Context, chatID int64, userID int) error
	UnbanChatMember(ctx context.Context, chatID int64, userID int) error
}

type service struct {
	apiEndpoint string
	apiToken    string
	Logger      log.Logger
}

func NewService(logger log.Logger, apiEndpoint string, apiToken string) *service {
	return &service{
		apiEndpoint: apiEndpoint,
		apiToken:    apiToken,
		Logger:      logger,
	}
}

func (s *service) buildURL(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", s.apiEndpoint, s.apiToken, method)
}

func (s *service) doRequest(ctx context.Context, url string, reqBody []byte) error {
	var tgResp APIResponse

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
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

func (s *service) SendMessage(ctx context.Context, chatID int64, text string) error {
	// https://core.telegram.org/bots/api#sendmessage
	const method = "sendMessage"

	req := SendMessageBody{
		ChatID:                chatID,
		Text:                  text,
		DisableWebPagePreview: true,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}

func (s *service) Reply(ctx context.Context, chatID int64, text string, messageID int) error {
	// https://core.telegram.org/bots/api#sendmessage
	const method = "sendMessage"

	req := SendMessageBody{
		ChatID:                chatID,
		Text:                  text,
		DisableWebPagePreview: true,
		ReplyToMessageID:      messageID,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}

func (s *service) KickChatMember(ctx context.Context, chatID int64, userID int) error {
	// https://core.telegram.org/bots/api#kickchatmember
	const method = "kickChatMember"

	req := KickChatMemberBody{
		ChatID: chatID,
		UserID: userID,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}

func (s *service) UnbanChatMember(ctx context.Context, chatID int64, userID int) error {
	// https://core.telegram.org/bots/api#unbanchatmember
	const method = "unbanChatMember"

	req := UnbanChatMemberBody{
		ChatID: chatID,
		UserID: userID,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}
