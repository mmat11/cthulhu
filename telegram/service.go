package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Service interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
	Reply(ctx context.Context, chatID int64, text string, messageID int) error
	KickChatMember(ctx context.Context, chatID int64, userID int) error
	UnbanChatMember(ctx context.Context, chatID int64, userID int) error
	DeleteMessage(ctx context.Context, chatID int64, messageID int) error
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		level.Error(s.Logger).Log("msg", "failed to create request", "err", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		level.Error(s.Logger).Log("msg", "failed reading response", "err", err)
	}

	json.Unmarshal(body, &tgResp)
	if !tgResp.Ok {
		return fmt.Errorf("api error: %s, code: %v", tgResp.Description, tgResp.ErrorCode)
	}
	return nil
}

func (s *service) SendMessage(ctx context.Context, chatID int64, text string) error {
	// https://core.telegram.org/bots/api#sendmessage
	const method = "sendMessage"

	reqJSON, err := json.Marshal(map[string]interface{}{
		"chat_id":                  chatID,
		"text":                     text,
		"disable_web_page_preview": true,
	})
	if err != nil {
		level.Error(s.Logger).Log("msg", "failed marshalling json", "err", err)
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}

func (s *service) Reply(ctx context.Context, chatID int64, text string, messageID int) error {
	// https://core.telegram.org/bots/api#sendmessage
	const method = "sendMessage"

	reqJSON, err := json.Marshal(map[string]interface{}{
		"chat_id":                  chatID,
		"text":                     text,
		"disable_web_page_preview": true,
		"reply_to_message_id":      messageID,
	})
	if err != nil {
		level.Error(s.Logger).Log("msg", "failed marshalling json", "err", err)
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}

func (s *service) KickChatMember(ctx context.Context, chatID int64, userID int) error {
	// https://core.telegram.org/bots/api#kickchatmember
	const method = "kickChatMember"

	reqJSON, err := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})
	if err != nil {
		level.Error(s.Logger).Log("msg", "failed marshalling json", "err", err)
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}

func (s *service) UnbanChatMember(ctx context.Context, chatID int64, userID int) error {
	// https://core.telegram.org/bots/api#unbanchatmember
	const method = "unbanChatMember"

	reqJSON, err := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})
	if err != nil {
		level.Error(s.Logger).Log("msg", "failed marshalling json", "err", err)
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}

func (s *service) DeleteMessage(ctx context.Context, chatID int64, messageID int) error {
	// https://core.telegram.org/bots/api#deletemessage
	const method = "deleteMessage"

	reqJSON, err := json.Marshal(map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	})
	if err != nil {
		level.Error(s.Logger).Log("msg", "failed marshalling json", "err", err)
		return err
	}

	return s.doRequest(ctx, s.buildURL(method), reqJSON)
}
