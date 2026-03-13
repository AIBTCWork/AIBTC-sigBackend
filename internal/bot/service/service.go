package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ecodeclub/ekit"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type BotServiceI interface {
	SendMessage(c context.Context, content string, ext ExtendFields) (e error)
	SendPhoto(c context.Context, photoUrl string, ext ExtendFields) (e error)
	SendVideo(c context.Context, ext ExtendFields) (e error)
	SendVoice(c context.Context, ext ExtendFields) (e error)
}

type botService struct {
	apiKey  string
	chatId  string
	baseUrl string
	client  *resty.Client
}

func NewBotService() BotServiceI {
	return &botService{
		apiKey:  viper.GetString("bot.api_key"),
		chatId:  viper.GetString("bot.chat_id"),
		baseUrl: viper.GetString("bot.base_url"),
		client:  resty.New(),
	}
}

func (s *botService) request(c context.Context, method messageType, payload map[string]string) error {
	url := fmt.Sprintf("%s%s/%s", s.baseUrl, s.apiKey, method)
	resp, err := s.client.R().
		SetContext(c).
		SetQueryParams(payload).
		Get(url)
	if err != nil {
		return err
	}
	fmt.Println(payload)
	if !resp.IsSuccess() {
		return fmt.Errorf("telegram error: %s", resp.String())
	}
	fmt.Println(resp.String())
	return nil
}

// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe

func (s *botService) SendMessage(c context.Context, content string, ext ExtendFields) (e error) {
	payload := map[string]string{
		"chat_id": s.chatId,
		"text":    content,
	}
	for k := range ext {
		payload[k], e = ext.Get(k).AsString()
	}
	e = s.request(c, text, payload)
	return
}

func (s *botService) SendPhoto(c context.Context, photoUrl string, ext ExtendFields) (e error) {
	payload := map[string]string{
		"chat_id": s.chatId,
		"photo":   photoUrl,
	}
	for k := range ext {
		payload[k], e = ext.Get(k).AsString()
	}
	e = s.request(c, photo, payload)
	return
}

func (s *botService) SendVideo(c context.Context, ext ExtendFields) (e error) {
	// TODO:
	return
}

func (s *botService) SendVoice(c context.Context, ext ExtendFields) (e error) {
	// TODO:
	return
}

type messageType string

const (
	text  messageType = "sendMessage"
	photo messageType = "sendPhoto"
	video messageType = "sendVideo"
	Voice messageType = "sendVoice"
)

type ExtendFields map[string]string

var errKeyNotFound = errors.New("没有找到对应的 key")

func (f ExtendFields) Get(key string) ekit.AnyValue {
	val, ok := f[key]
	if !ok {
		return ekit.AnyValue{
			Err: fmt.Errorf("%w, key %s", errKeyNotFound, key),
		}
	}
	return ekit.AnyValue{Val: val}
}
