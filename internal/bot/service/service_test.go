package service

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

func TestSender(t *testing.T) {
	// 初始化 viper 配置
	viper.SetConfigFile("../../../conf/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("读取配置文件失败: %v", err)
	}

	sender := NewBotService()
	e := sender.SendMessage(context.Background(), "hello world", ExtendFields{})
	t.Log(e)

}

// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
func TestSenderAuth(t *testing.T) {
	url := "http://api.telegram.org/bot8691164628:AAFx41THgAd9zB56OtnJ-pXEWiezlLRCW3U/getMe"
	c := resty.New()
	resp, e := c.R().Get(url)
	if e != nil {
		t.Error(e)
	}

	t.Log(resp.String())
}
