package app

import (
	"github.com/atlassian/go-sentry-api"
	"github.com/evgeny-klyopov/sentry-notifier/app/message"
	"github.com/evgeny-klyopov/sentry-notifier/app/sender"
	"github.com/evgeny-klyopov/sentry-notifier/config"
	"strconv"
)

func (a *App) getSentryProjects(cfgOrg config.Organization, client *sentry.Client) (*[]sentry.Project, error) {
	var allProject bool
	if len(*cfgOrg.Projects) == 0 {
		allProject = true
	}

	sentryProjects, _, err := client.GetProjects()
	if err != nil {
		return nil, err
	}

	var projects []sentry.Project
	for _, project := range sentryProjects {
		if allProject == true {
			projects = append(projects, project)
		} else {
			for _, configProject := range *cfgOrg.Projects {
				if configProject.Name == project.Name {
					projects = append(projects, project)
				}
			}
		}
	}

	return &projects, nil
}

func (a *App) sendByProject(project sentry.Project, issues []sentry.Issue) error {
	if len(issues) == 0 {
		return nil
	}

	suffixFileLog := *project.Organization.ID + "_" + project.ID + ".log"

	msgObject, err := message.CreateObject(message.MdFormatType)
	if err != nil {
		return err
	}

	if len(*a.config.Default.Notifications.Telegram) > 0 {
		prefixFileLog := a.dirLog + "/telegram"
		st, err := sender.CreateObject(sender.TelegramType)
		if err != nil {
			return err
		}

		for _, cfg := range *a.config.Default.Notifications.Telegram {
			err = a.processSend(
				project,
				issues,
				msgObject,
				st,
				sender.Config{Telegram: cfg},
				prefixFileLog+strconv.FormatInt(cfg.ChatId, 10)+"_"+suffixFileLog,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
