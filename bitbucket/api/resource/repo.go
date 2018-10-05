package resource

import "github.com/auvn/go-atlassian/bitbucket/api"

type Repo struct {
	Parent api.Resource
	Slug   string
}

func (r Repo) Path() api.Path {
	return api.Path{
		Path: r.Parent.Path().Path + "/repos/" + r.Slug,
	}
}

func (r Repo) PullRequest(id int64) PullRequest {
	return PullRequest{
		Parent: r,
		ID:     id,
	}
}
