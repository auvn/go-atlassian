package api

import (
	"strconv"
	"strings"
)

type Resource interface {
	Path() Path
}

type Path struct {
	Path   string
	Params PathParams
}

func (p Path) String() string {
	return p.Path + p.Params.Join()
}

type PathParams map[string]string

func (p PathParams) WithString(key, value string) PathParams {
	if p == nil {
		p = PathParams{}
	}

	if value == "" {
		return p
	}

	p[key] = value

	return p
}

func (p PathParams) WithInt64(key string, value int64) PathParams {
	return p.WithString(key, strconv.Itoa(int(value)))
}

func (p PathParams) Join() string {
	paramsLen := len(p)
	if paramsLen > 0 {
		params := make([]string, 0, paramsLen)
		for k, v := range p {
			params = append(params, k+"="+v)
		}
		return "?" + strings.Join(params, "&")
	}
	return ""
}
