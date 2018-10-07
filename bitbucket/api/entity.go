package api

type Link struct {
	Href string `json:"href"`
}

type Links struct {
	Self []Link `json:"self"`
}

type User struct {
	EmailAddress string `json:"emailAddress"`
}

type Project struct {
	Key string `json:"key"`
}

type Repository struct {
	Slug    string  `json:"slug"`
	Project Project `json:"project"`
}

type Ref struct {
	Repository Repository `json:"repository"`
}

type PullRequestAuthor struct {
	User User `json:"user"`
}

type PullRequest struct {
	Links       Links             `json:"links"`
	ID          int64             `json:"id"`
	Version     int64             `json:"version"`
	Title       string            `json:"title"`
	FromRef     Ref               `json:"fromRef"`
	ToRef       Ref               `json:"toRef"`
	Author      PullRequestAuthor `json:"author"`
	CreatedDate int64             `json:"createdDate"`
	UpdatedDate int64             `json:"updatedDate"`
}

type PullRequests struct {
	Page
	Values []PullRequest
}
