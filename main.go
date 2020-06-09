package main

import (
	"sentry-notifier/app"
	"sentry-notifier/config"
)

func main()  {
	cfg, err := config.GetConfig("config.json")

	if err != nil {
		panic(err)
	}

	err = app.NewApp(*cfg).Run()
	if err != nil {
		panic(err)
	}
}