package activity

import (
	"context"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
	"github.com/auvn/go-json/jsonutil"
)

type activitiesPage struct {
	api.Page
	Values []jsonutil.Object `json:"values"`
}

func GetPage(ctx context.Context, g *atlassian.RestClient, path api.Path) (*Activities, error) {
	var resp activitiesPage
	if err := g.Get(ctx, path.String(), &resp); err != nil {
		return nil, err
	}

	return &Activities{
		NextPage: resp.Page.Next(path),
		IsLast:   resp.IsLastPage,
		Values:   resp.Values,
	}, nil
}
