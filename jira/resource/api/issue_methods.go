package api

import (
	"context"

	"github.com/auvn/go-juno/juno/git/bitbucket/resource/resourceutil"
)

func (tt IssueTransitions) Post(ctx context.Context, p resourceutil.Poster) error {
	return p.Post(ctx, tt.URL(), nil)
}
