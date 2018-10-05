package resource

import (
	"strconv"

	"github.com/auvn/go-atlassian/bitbucket/api"
)

type PullRequest struct {
	Parent api.Resource
	ID     int64
}

func (p PullRequest) Path() api.Path {
	return api.Path{
		Path: p.Parent.Path().Path + "/pull-requests/" + strconv.Itoa(int(p.ID)),
	}
}

func (p PullRequest) Activities() PullRequestActivities {
	return PullRequestActivities{
		Parent: p,
	}
}

type PullRequestActivities struct {
	Parent api.Resource
}

func (a PullRequestActivities) Path() api.Path {
	return api.Path{
		Path: a.Parent.Path().Path + "/activities",
	}
}
