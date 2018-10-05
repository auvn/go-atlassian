package activity

import (
	"context"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
)

func Get(ctx context.Context, g *atlassian.RestClient, path api.Resource) (*GetActivitiesResponse, error) {
	var resp GetActivitiesResponse
	if err := g.Get(ctx, path.Path().String(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
