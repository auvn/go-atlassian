package dashboard

import (
	"strconv"

	"github.com/auvn/go-atlassian/bitbucket/api"
)

type PullRequestsParams struct {
	PullRequests PullRequests

	State string
	Role  string
	Order string
	Limit string
}

func (p PullRequestsParams) WithState() PullRequestsStateParams {
	return PullRequestsStateParams(p)
}

func (p PullRequestsParams) WithRole() PullRequestsRoleParams {
	return PullRequestsRoleParams(p)
}

func (p PullRequestsParams) WithOrder() PullRequestsOrderParams {
	return PullRequestsOrderParams(p)
}

func (p PullRequestsParams) WithLimit(limit int) PullRequestsParams {
	p.Limit = strconv.Itoa(limit)
	return p
}

func (p PullRequestsParams) Path() api.Path {
	return api.Path{
		Path: p.PullRequests.Path().Path,
		Params: api.PathParams{}.
			WithString("state", p.State).
			WithString("role", p.Role).
			WithString("order", p.Order).
			WithString("limit", p.Limit),
	}
}

type PullRequestsStateParams PullRequestsParams

func (s PullRequestsStateParams) Open() PullRequestsParams {
	s.State = "OPEN"
	return PullRequestsParams(s)
}

func (s PullRequestsStateParams) Declined() PullRequestsParams {
	s.State = "DECLINED"
	return PullRequestsParams(s)
}

func (s PullRequestsStateParams) Merged() PullRequestsParams {
	s.State = "MERGED"
	return PullRequestsParams(s)
}

type PullRequestsRoleParams PullRequestsParams

func (r PullRequestsRoleParams) Reviewer() PullRequestsParams {
	r.Role = "REVIEWER"
	return PullRequestsParams(r)
}

func (r PullRequestsRoleParams) Author() PullRequestsParams {
	r.Role = "AUTHOR"
	return PullRequestsParams(r)
}

func (r PullRequestsRoleParams) Participant() PullRequestsParams {
	r.Role = "PARTICIPANT"
	return PullRequestsParams(r)
}

type PullRequestsOrderParams PullRequestsParams

func (r PullRequestsOrderParams) Newest() PullRequestsParams {
	r.Order = "NEWEST"
	return PullRequestsParams(r)
}

func (r PullRequestsOrderParams) Oldest() PullRequestsParams {
	r.Order = "OLDEST"
	return PullRequestsParams(r)
}
