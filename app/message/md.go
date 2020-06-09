package message

import (
	"github.com/atlassian/go-sentry-api"
	"strconv"
	"strings"
)

const MdFormatType = "md"

type MD struct {

}

func (m MD) Build(project sentry.Project, issue sentry.Issue) string {
	message := []string{
		"*" + project.Organization.Name + " [" + *project.Organization.ID + "]*",
		"-------------------------------------",
		"*Project:* " + project.Name,
		"*Project ID:* " + project.ID,
		"",
		"*Issue ID:* " + *issue.ID,
		"*Title:* " + *issue.Title,
		"*Type:* " + *issue.Type,
		"*FirstSeen:* " + issue.FirstSeen.Format("02.01.2006 15:04:05"),
		"*LastSeen:* " + issue.LastSeen.Format("02.01.2006 15:04:05"),
		"*Status:* " + string(*issue.Status),
		"*Short ID:* " + *issue.ShortID,
		"*Level:* " + *issue.Level,
		"*Culprit:* " + *issue.Culprit,
		"*UserCount:* " + strconv.Itoa(*issue.UserCount),
		"*Permalink:* " + *issue.Permalink,
		"",
	}

	return strings.Join(message, "\n")
}