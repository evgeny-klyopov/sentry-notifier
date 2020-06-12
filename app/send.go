package app

import (
	"github.com/atlassian/go-sentry-api"
	"github.com/evgeny-klyopov/sentry-notifier/app/message"
	"github.com/evgeny-klyopov/sentry-notifier/app/sender"
	"os"
	"sort"
	"strconv"
)

func (a *App) processSend(project sentry.Project, issues []sentry.Issue, msgObject message.Messenger, st sender.Sender, senderConfig sender.Config, logFile string) error {
	err := st.SetClient(senderConfig)
	if err != nil {
		return err
	}

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
