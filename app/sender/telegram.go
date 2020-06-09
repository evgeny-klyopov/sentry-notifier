package sender

import (
	"github.com/evgeny-klyopov/telegram-simple-message"
)

const TelegramType = "telegram"

type Telegram struct {
	client *telegram.Client
}

func (t *Telegram) SetClient(cfg Config) error {
	var err error
	t.client, err = telegram.NewClient(cfg.Telegram)

	if err != nil {
		return err
	}

	return nil
}

func (t *Telegram) Send(message string) error {
	return t.client.Send(message, telegram.MarkdownTypeMessage)
}
