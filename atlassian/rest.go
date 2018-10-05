package atlassian

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type RestClient struct {
	Client Client
}

func (c *RestClient) do(ctx context.Context, method, path string, dest interface{}) error {
	req, err := c.Client.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return resp.UnmarshalJSON(dest)
}

func (c *RestClient) Get(ctx context.Context, path string, dest interface{}) error {
	if err := c.do(ctx, http.MethodGet, path, dest); err != nil {
		return errors.Wrap(err, "failed to perform get request")
	}
	return nil
}

func (c *RestClient) Post(ctx context.Context, path string, dest interface{}) error {
	if err := c.do(ctx, http.MethodPost, path, dest); err != nil {
		return errors.Wrap(err, "failed to perform post request")
	}
	return nil
}
