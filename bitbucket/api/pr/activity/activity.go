package activity

import (
	"github.com/auvn/go-atlassian/bitbucket/api"
	"github.com/auvn/go-atlassian/bitbucket/api/pr"
	"github.com/auvn/go-json/jsonutil"
)

type GetActivitiesResponse struct {
	api.Page
	Values []jsonutil.Object `json:"values"`
}

const (
	ActionCommented = "COMMENTED"
)

func ActionOf(obj jsonutil.Object) string {
	const keyAction = "action"

	action, _ := obj.Value(keyAction)
	return action.String()
}

type Comment struct {
	Comment       pr.Comment    `json:"comment"`
	CommentAnchor CommentAnchor `json:"commentAnchor"`
}

type CommentAnchor struct {
	Line int    `json:"line"`
	Path string `json:"path"`
}
