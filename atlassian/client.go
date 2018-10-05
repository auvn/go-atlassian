package atlassian

import (
	"bytes"
	"context"
	"encoding/json"
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
