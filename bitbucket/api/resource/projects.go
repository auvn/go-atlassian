package resource

import "github.com/auvn/go-atlassian/bitbucket/api"

type Project struct {
	Parent api.Resource
	Key    string
}

func (p Project) Path() api.Path {
	return api.Path{
		Path: p.Parent.Path().Path + "/projects/" + p.Key,
	}
}

func (p Project) Repo(slug string) Repo {
	return Repo{
		Parent: p,
		Slug:   slug,
	}
}
