package app

import (
	"github.com/atlassian/go-sentry-api"
	"github.com/evgeny-klyopov/sentry-notifier/config"
	"os"
	"os/signal"
	"time"
)

func (a *App) process() error {
	var err error

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	stop := make(chan bool)
	errorChannel := make(chan error)

	for _, cfgOrg := range a.config.Organization {
		cfgOrg := cfgOrg
		go func() {
			defer func() { stop <- true }()
			for {
				select {
				case <-stop:
					return
				default:
					errorChannel <- a.handler(cfgOrg)
				}
			}
		}()
	}

	go func() {
		for {
			err = <-errorChannel
			if err != nil {
				signals <- os.Interrupt
			}
		}
	}()

	// ctrl c
	<-signals
	stop <- true
	<-stop

	return err
}

func (a *App) handler(cfgOrg config.Organization) error {
	client, _ := sentry.NewClient(cfgOrg.Token, nil, nil)
	projects, err := a.getSentryProjects(cfgOrg, client)
	if err != nil {
		return err
	}

	org, err := client.GetOrganization(cfgOrg.Name)
	if err != nil {
		return err
	}
	for _, project := range *projects {
		issues, _, _ := client.GetIssues(org, project, &a.config.Default.Sentry.IssueFilter.StatsPeriod, nil, &a.config.Default.Sentry.IssueFilter.Query)
		err := a.sendByProject(project, issues)
		if err != nil {
			return err
		}
	}

	time.Sleep(time.Duration(a.config.Default.Sentry.WaitTime) * time.Second)

	return nil
}
