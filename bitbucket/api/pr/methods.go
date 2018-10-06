package pr

import (
	"context"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
)

type pullRequestsPage struct {
	api.Page
	Values []api.PullRequest `json:"values"`
}

func GetPage(ctx context.Context, g *atlassian.RestClient, path api.Path) (*api.PullRequests, error) {
	var resp pullRequestsPage
	if err := g.Get(ctx, path.String(), &resp); err != nil {
		return nil, err
	}
	return &api.PullRequests{
		NextPage: resp.Page.Next(path),
		IsLast:   resp.Page.IsLastPage,
		Values:   resp.Values,
	}, nil
}
