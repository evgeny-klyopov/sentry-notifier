package config

import "github.com/evgeny-klyopov/telegram-simple-message"

type notifications struct {
	Telegram *[]telegram.Config `json:"telegram"`
}
type issueFilter struct {
	Query       string `json:"query"`
	StatsPeriod string `json:"stats-period"`
}

type sentry struct {
	IssueFilter issueFilter `json:"issue_filter"`
	WaitTime    int64       `json:"wait_time"`
}
type Setting struct {
	Sentry        sentry        `json:"sentry"`
	Notifications notifications `json:"notifications"`
}

type project struct {
	Name    string   `json:"name"`
	Setting *Setting `json:"setting"`
}
type Organization struct {
	Name     string     `json:"name"  validate:"required"`
	Token    string     `json:"token"  validate:"required"`
	Projects *[]project `json:"projects"`
	Setting  *Setting   `json:"setting"`
}

type Config struct {
	Organization []Organization `json:"organization" validate:"required,dive,required"`
	Default      Setting        `json:"default" validate:"required,dive,required"`
}
