package jira

import "github.com/auvn/go-atlassian/atlassian"

type Rest struct {
	Client  atlassian.Client
	BaseURL string
}
