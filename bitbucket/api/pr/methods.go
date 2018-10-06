package pr

import (
	"context"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
)

func GetPage(ctx context.Context, g *atlassian.RestClient, path api.Path) (*api.PullRequests, error) {
	var resp api.PullRequests
	if err := g.Get(ctx, path.String(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
