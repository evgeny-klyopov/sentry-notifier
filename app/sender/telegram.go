package sender

import (
	"context"
	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
	"net"
	"net/http"
	"sentry-notifier/config"
	"strings"
)

const TelegramType = "telegram"

type Telegram struct {
	config config.Telegram
}

func (t *Telegram) SetConfig(cfg Config) Sender{
	t.config = cfg.Telegram
	return t
}

func (t *Telegram) Send(message string) error{
	var bot *tgbotapi.BotAPI

	if t.config.UseProxy {
		setting := strings.Split(t.config.Proxy, "@")
		authData := strings.Split(setting[0], ":")

		dialer, err := proxy.SOCKS5(
			"tcp",
			setting[1],
			&proxy.Auth{User: authData[0], Password:  authData[1]},
			proxy.Direct,
		)

		if err != nil {
			return err
		}

		client := &http.Client{Transport: &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}}}

		bot, err = tgbotapi.NewBotAPIWithClient(t.config.Token, client)
		if err != nil {
			return err
		}
	} else {
		var err error
		bot, err = tgbotapi.NewBotAPI(t.config.Token)
		if err != nil {
			return err
		}
	}

	msg := tgbotapi.NewMessage(t.config.ChatId, message)
	msg.ParseMode = "markdown"
	_, err := bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}