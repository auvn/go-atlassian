package pr

import (
	"context"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
)

func GetPage(ctx context.Context, g *atlassian.RestClient, path api.Resource) (*api.PullRequestsPage, error) {
	var resp api.PullRequestsPage
	if err := g.Get(ctx, path.Path().String(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
