package practivity

import (
	"context"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
	"github.com/auvn/go-atlassian/bitbucket/api/pr"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity/activityconv"
	"github.com/auvn/go-atlassian/bitbucket/api/resource"
	"github.com/auvn/go-atlassian/bitbucket/api/resource/dashboard"
	"github.com/auvn/go-atlassian/bitbucketutil"
	"github.com/auvn/go-shell/output"
	"github.com/auvn/go-shell/strfmt"
	"golang.org/x/sync/errgroup"
)

var bitbucketAPI = resource.Latest{}

func List(client *atlassian.DefaultClient, maxAge time.Duration) (*PullRequests, error) {
	prs, err := bitbucketutil.GetPullRequests(client,
		dashboard.New(bitbucketAPI).
			PullRequests().
			WithParams().
			WithRole().Author().
			WithState().Open())
	if err != nil {
		return nil, err
	}

	var pullRequests PullRequests

	if len(prs) == 0 {
		return &pullRequests, nil
	}

	ch, err := listPullRequests(client, time.Now().UTC().Add(-maxAge), prs)
	if err != nil {
		return nil, err
	}

	for prRef := range ch {
		pullRequests.PRs = append(pullRequests.PRs, *prRef)
	}

	sort.Slice(pullRequests.PRs, func(i, j int) bool {
		return pullRequests.PRs[i].Order > pullRequests.PRs[j].Order
	})

	return &pullRequests, nil
}

func listPullRequests(client *atlassian.DefaultClient, tm time.Time, prs []api.PullRequest) (<-chan *PullRequest, error) {
	pullRequestsChan := make(chan *PullRequest, len(prs))
	defer close(pullRequestsChan)

	group, _ := errgroup.WithContext(context.TODO())
	for i := range prs {
		i := i
		group.Go(func() error {
			prRef, err := listPullRequest(client, tm, prs[i])
			if err != nil {
				return err
			}

			pullRequestsChan <- prRef

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return pullRequestsChan, nil
}

func listPullRequest(client *atlassian.DefaultClient, tm time.Time, pr api.PullRequest) (*PullRequest, error) {
	repo := pr.FromRef.Repository
	activities, err := bitbucketutil.GetActivities(client,
		bitbucketAPI.
			Project(repo.Project.Key).
			Repo(repo.Slug).
			PullRequest(pr.ID).
			Activities())
	if err != nil {
		return nil, err
	}

	prRef := PullRequest{
		Title: pr.Title,
		Link:  pr.Links.Self[0].Href,
		Order: pr.UpdatedDate,
	}

	if len(activities) == 0 {
		return &prRef, nil
	}

	for j := len(activities) - 1; j >= 0; j-- {
		if activity.ActionOf(activities[j]) != activity.ActionCommented {
			continue
		}

		commentActivity := activityconv.CommentFromObject(activities[j])

		latestComment := bitbucketutil.LatestComment(commentActivity.Comment)
		if latestComment.
			UpdatedAt().Before(tm) {
			continue
		}

		prRef.Activities = append(prRef.Activities, Comment{
			Comment:   commentActivity,
			Highlight: latestComment.Author.EmailAddress != pr.Author.User.EmailAddress,
		})
	}
	return &prRef, nil
}

type PullRequests struct {
	PRs []PullRequest
}

func (prs PullRequests) Fprintf(w io.Writer) {
	for i := range prs.PRs {
		prs.PRs[i].Fprint(w)
		fmt.Fprintln(w)
	}

}

type PullRequest struct {
	Title      string
	Link       string
	Order      int64
	Activities []Comment
	Highlight  bool
}

func (pr PullRequest) Fprint(w io.Writer) {
	strfmt.Fprintf(w, strfmt.StyleBold, "%s\n", pr.Title)
	fmt.Fprintf(w, "%s\n\n", pr.Link)

	for i := range pr.Activities {
		pr.Activities[i].Fprint(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
	}
}

type Comment struct {
	activity.Comment
	Highlight bool
}

func (a Comment) Fprint(w io.Writer) {
	if a.Highlight {
		w = output.NewPrefixWriter(w, strfmt.Bold("  =/=/= "))
	} else {
		w = output.NewPrefixWriter(w, "        ")
	}

	if a.CommentAnchor.Path != "" {
		strfmt.Fprintf(w, strfmt.StyleBold,
			"Line: %d Path: %s:\n", a.CommentAnchor.Line, a.CommentAnchor.Path)
	}
	PRComments{Comments: []pr.Comment{a.Comment.Comment}}.Fprint(w)
}

type PRComments struct {
	indentWidth int
	Comments    []pr.Comment
}

func (cc PRComments) Fprint(w io.Writer) {
	indentWriter := output.NewSameIndentWriter(w, cc.indentWidth)
	for _, c := range cc.Comments {
		fmt.Fprintf(indentWriter, "%s: %s\n", strfmt.Bold(c.Author.EmailAddress), c.Text)

		PRComments{indentWidth: cc.indentWidth + 2, Comments: c.Comments}.Fprint(w)
	}
}
