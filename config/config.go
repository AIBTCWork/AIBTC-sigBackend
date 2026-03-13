package config

import (
	"time"
)

type Config struct {
	Server ServerConfig `json:"server" yaml:"server"` // server config
	Bot    BotConfig    `json:"bot" yaml:"bot"`
}
type ServerConfig struct {
	Domain           string          `json:"domain" yaml:"domain"`                       // domain
	Port             uint            `json:"port" yaml:"port"`                           // web server port
	ReadTimeout      time.Duration   `json:"read_timeout" yaml:"read_timeout"`           // read timeout
	WriteTimeout     time.Duration   `json:"write_timeout" yaml:"write_timeout"`         // write timeout
	GracefulShutdown time.Duration   `json:"graceful_shutdown" yaml:"graceful_shutdown"` // graceful shutdown time
	Mode             string          `json:"mode" yaml:"mode"`                           // environment(local/dev/pre/prod)
	Whitelist        map[string]bool `json:"whitelist" yaml:"whitelist"`                 // whitelist
}

type BotConfig struct {
	ApiKey  string `json:"api_key" yaml:"api_key"` // telegram bot api key
	ChatId  string `json:"chat_id" yaml:"chat_id"` // telegram chat id
	BaseUrl string `json:"base_url" yaml:"base_url"`
}
