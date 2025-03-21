package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DomainName  string `mapstructure:"DOMAIN_NAME"`
	Production  bool   `mapstructure:"PRODUCTION"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`

	ClientOrigin   string
	ServerOrigin   string
	ServerPort     string `mapstructure:"SERVER_PORT"`
	AllowedOrigins []string
	CookieDomain   string

	TelegramBotToken string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID   string `mapstructure:"TELEGRAM_CHAT_ID"`
}

var Env *Config

func LoadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var protocol, cookieDomain, clientOrigin, serverOrigin string
	if viper.GetBool("PRODUCTION") {
		protocol = "https://"
		cookieDomain = "." + viper.GetString("DOMAIN_NAME")
		clientOrigin = protocol + viper.GetString("DOMAIN_NAME")
		serverOrigin = protocol + "api." + viper.GetString("DOMAIN_NAME")
	} else if !viper.GetBool("PRODUCTION") {
		protocol = "http://"
		cookieDomain = viper.GetString("DOMAIN_NAME")
		clientOrigin = protocol + viper.GetString("DOMAIN_NAME") + ":" + viper.GetString("CLIENT_PORT")
		serverOrigin = protocol + viper.GetString("DOMAIN_NAME") + ":" + viper.GetString("SERVER_PORT")

	}
	allowedOrigins := []string{clientOrigin, serverOrigin}

	Env = &Config{
		DomainName:     viper.GetString("DOMAIN_NAME"),
		Production:     viper.GetBool("PRODUCTION"),
		DatabaseUrl:    viper.GetString("DATABASE_URL"),
		ClientOrigin:   clientOrigin,
		ServerOrigin:   serverOrigin,
		ServerPort:     viper.GetString("SERVER_PORT"),
		AllowedOrigins: allowedOrigins,
		CookieDomain:   cookieDomain,

		TelegramBotToken: viper.GetString("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   viper.GetString("TELEGRAM_CHAT_ID"),
	}

	return nil
}
