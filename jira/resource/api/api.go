package api

import "github.com/auvn/go-juno/juno/git/bitbucket/resource"

type API struct{}

func (a API) Issue(id string) Issue {
	return Issue{
		API: a,
		ID:  id,
	}
}

func (a API) URL() resource.URL {
	return resource.URL{
		Path: "/api/latest",
	}
}
