package sender

import (
	"github.com/evgeny-klyopov/telegram-simple-message"
)

type Config struct {
	Telegram telegram.Config
}

type Sender interface {
	Send(message string) error
	SetClient(cfg Config) error
}
