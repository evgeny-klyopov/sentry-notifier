package message

import "github.com/atlassian/go-sentry-api"

type Messenger interface {
	Build(project sentry.Project, issues sentry.Issue) string
}
