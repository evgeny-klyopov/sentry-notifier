package app

import (
	"github.com/atlassian/go-sentry-api"
	"io"
	"io/ioutil"
	"os"
	"sentry-notifier/app/message"
	"sentry-notifier/app/sender"
	"sentry-notifier/config"
	"sort"
	"strconv"
	"time"
)

type sendMessage struct{
	message string
	stringId string
	id int64
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

	for _, cfgOrg := range a.config.Organization {
		for {
			err := a.checkAndSend(cfgOrg)
			if err != nil {
				return err
			}

			time.Sleep(time.Duration(a.config.Default.Sentry.WaitTime) * time.Second)
		}
	}

	return nil
}


func (a *App) checkAndSend(cfgOrg config.Organization) error {
	client, _:= sentry.NewClient(cfgOrg.Token, nil, nil)
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
		err := a.sendNotification(project, issues)
		if err != nil {
			return err
		}
	}

	return nil
}





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
			projects = append(projects,  project)
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

func (a *App) sendNotification(project sentry.Project, issues []sentry.Issue) error {
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
				st.SetConfig(sender.Config{Telegram: cfg}),
				prefixFileLog+strconv.FormatInt(cfg.ChatId, 10)+"_"+suffixFileLog,
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (a *App) getMessages(project sentry.Project, issues []sentry.Issue, logFile string, msgObject message.Messenger) (*[]sendMessage, error) {
	var lastIssueID *string
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		empty := "0"
		lastIssueID = &empty
	} else {
		var err error
		lastIssueID, err = a.readLog(logFile)
		if err != nil {
			return nil, err
		}
	}

	lastIssueIDInt, _ := strconv.ParseInt(*lastIssueID, 10, 64)

	var messages []sendMessage
	for _, issue := range issues {
		id, _ := strconv.ParseInt(*issue.ID, 10, 64)

		if id < lastIssueIDInt {
			continue
		}

		if id > lastIssueIDInt {
			messages = append(messages, sendMessage{
				message:  msgObject.Build(project, issue),
				id:       id,
				stringId: *issue.ID,
			})
		}
	}

	return &messages, nil
}
func (a *App) readLog(path string) (*string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := string(data)
	return &result, nil
}


func (a *App) writeToLog(path string, lastIssueID *string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()


	_, err = io.WriteString(file, *lastIssueID)
	if err != nil {
		return err
	}
	return file.Sync()
}

func (a *App) sendMessages(messages *[]sendMessage, st sender.Sender, logFile string) error {
	if len(*messages) > 0 {
		sort.SliceStable(*messages, func(i, j int) bool {
			return (*messages)[i].id < (*messages)[j].id
		})

		var lastIssueID *string
		for _, msg := range *messages {
			err := st.Send(msg.message)
			if err != nil {
				_ = a.writeToLog(logFile, lastIssueID)
				return err
			}
			lastIssueID = &msg.stringId
		}

		err := a.writeToLog(logFile, lastIssueID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) processSend(project sentry.Project, issues []sentry.Issue, msgObject message.Messenger, st sender.Sender, logFile string) error {
	messages, err := a.getMessages(
		project,
		issues,
		logFile,
		msgObject,
	)

	if err != nil {
		return err
	}

	return a.sendMessages(messages, st, logFile)
}

























