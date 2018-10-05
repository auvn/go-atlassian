package resource

import "github.com/auvn/go-atlassian/bitbucket/api"

type Latest struct{}

func (l Latest) Path() api.Path {
	return api.Path{
		Path: "/api/latest",
	}
}

func (l Latest) Project(key string) Project {
	return Project{
		Parent: l,
		Key:    key,
	}
}
