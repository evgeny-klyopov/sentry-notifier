package sender

import "sentry-notifier/config"

type Config struct {
	Telegram config.Telegram
}

type Sender interface {
	Send(message string) error
	SetConfig(cfg Config) Sender
}