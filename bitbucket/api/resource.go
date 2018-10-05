package api

import "strings"

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

type PathParams []string

func (p PathParams) WithString(key, value string) PathParams {
	if value == "" {
		return p
	}
	return append(p, key+"="+value)
}

func (p PathParams) Join() string {
	if len(p) > 0 {
		return "?" + strings.Join(p, "&")
	}
	return ""
}
