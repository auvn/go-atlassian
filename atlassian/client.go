package atlassian

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Response struct {
	*http.Response
}

func (r *Response) UnmarshalJSON(dest interface{}) error {
	if dest == nil {
		return nil
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dest); err != nil {
		return errors.Wrap(err, "failed to decode body")
	}
	return nil
}

type Client struct {
	Http    http.Client
	Auth    string
	BaseURL string
}

func marshalBody(body interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	if body == nil {
		return &buf, nil
	}

	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(body); err != nil {
		return nil, errors.Wrap(err, "failed to encode body")
	}
	return &buf, nil
}

func (c *Client) setDefaultHeaders(req *http.Request) {
	if c.Auth != "" {
		req.Header.Set("Authorization", c.Auth)
	}
	req.Header.Set("X-Atlassian-Token", "no-check")
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	bodyReader, err := marshalBody(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	c.setDefaultHeaders(req)

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.Http.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to perform http request")
	}

	return &Response{resp}, nil
}

type DefaultClient struct {
	Client Client
}

func (c *DefaultClient) do(ctx context.Context, method, path string, dest interface{}) error {
	req, err := c.Client.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	//TODO
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
		case http.StatusUnauthorized:
		case http.StatusForbidden:
		case http.StatusNotFound:
		case http.StatusMethodNotAllowed:
		case http.StatusConflict:
		case http.StatusUnsupportedMediaType:
		case http.StatusInternalServerError:
		}

		bb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("%s: %s", resp.Status, (bb))
	}

	return resp.UnmarshalJSON(dest)
}

func (c *DefaultClient) Get(ctx context.Context, path string, dest interface{}) error {
	if err := c.do(ctx, http.MethodGet, path, dest); err != nil {
		return errors.Wrap(err, "failed to perform get request")
	}
	return nil
}

func (c *DefaultClient) Post(ctx context.Context, path string, dest interface{}) error {
	if err := c.do(ctx, http.MethodPost, path, dest); err != nil {
		return errors.Wrap(err, "failed to perform post request")
	}
	return nil
}
