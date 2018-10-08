package api

import "github.com/auvn/go-juno/juno/git/bitbucket/resource"

type Issue struct {
	API
	ID string
}

func (i Issue) Transitions() IssueTransitions {
	return IssueTransitions{
		Issue: i,
	}
}

func (i Issue) URL() resource.URL {
	return resource.URL{
		Path: i.API.URL().Path + "/issue/" + i.ID,
	}
}

type IssueTransitions struct {
	Issue
}

func (tt IssueTransitions) URL() resource.URL {
	return resource.URL{
		Path: tt.Issue.URL().Path + "/transitions",
	}
}
