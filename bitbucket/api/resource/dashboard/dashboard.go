package dashboard

import "github.com/auvn/go-atlassian/bitbucket/api"

type Dashboard struct {
	Parent api.Resource
}

func (d Dashboard) PullRequests() PullRequests {
	return PullRequests{
		Parent: d,
	}
}

func (d Dashboard) Path() api.Path {
	return api.Path{
		Path: d.Parent.Path().Path + "/dashboard",
	}
}

func New(parent api.Resource) Dashboard {
	return Dashboard{parent}
}

type PullRequests struct {
	Parent api.Resource
}

func (prs PullRequests) WithParams() PullRequestsParams {
	return PullRequestsParams{
		PullRequests: prs,
	}
}

func (prs PullRequests) Path() api.Path {
	return api.Path{
		Path: prs.Parent.Path().Path + "/pull-requests",
	}
}
