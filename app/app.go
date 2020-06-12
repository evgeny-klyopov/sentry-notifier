package app

import (
	"github.com/evgeny-klyopov/sentry-notifier/config"
	"os"
)

type sendMessage struct {
	message  string
	stringId string
	id       int64
}

type App struct {
	config config.Config
	dirLog string
}

func NewApp(cfg config.Config) *App {
	return &App{
		config: cfg,
		dirLog: "logs",
	}
}

func (a *App) Run() error {
	if _, err := os.Stat(a.dirLog); os.IsNotExist(err) {
		err = os.Mkdir(a.dirLog, 0750)
		if err != nil {
			return err
		}
	}

	return a.process()
}
