package activity

import (
	"context"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
)

func GetPage(ctx context.Context, g *atlassian.DefaultClient, path api.Path) (*Activities, error) {
	var resp Activities
	if err := g.Get(ctx, path.String(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
