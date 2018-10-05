package pr

import (
	"time"

	"github.com/auvn/go-atlassian/bitbucket/api"
)

type Comment struct {
	Author      api.User  `json:"author"`
	Text        string    `json:"text"`
	Comments    []Comment `json:"comments"`
	UpdatedDate int64     `json:"updatedDate"`
}

func (c Comment) UpdatedAt() time.Time {
	return time.Unix(c.UpdatedDate, 0)
}
