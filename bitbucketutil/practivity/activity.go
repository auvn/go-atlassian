package practivity

import (
	"context"
	"fmt"
	"io"
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

type ListParams struct {
	MaxAge   time.Duration
	IsAuthor bool
}

func List(client *atlassian.DefaultClient, params ListParams) (*PullRequests, error) {

	resourceParams := dashboard.New(bitbucketAPI).
		PullRequests().
		WithParams().
		WithState().Open()
	if params.IsAuthor {
		resourceParams = resourceParams.WithRole().Author()
	}

	prs, err := bitbucketutil.GetPullRequests(client, resourceParams)
	if err != nil {
		return nil, err
	}

	var pullRequests PullRequests

	if len(prs) == 0 {
		return &pullRequests, nil
	}

	var commentsAfter time.Time
	if params.MaxAge > 0 {
		commentsAfter = time.Now().UTC().Add(-params.MaxAge)
	}

	pullRequests.PRs, err = listPullRequests(client, commentsAfter, prs)
	if err != nil {
		return nil, err
	}

	return &pullRequests, nil
}

func listPullRequests(client *atlassian.DefaultClient, tm time.Time, prs []api.PullRequest) ([]PullRequest, error) {
	pullRequests := make([]PullRequest, len(prs))

	group, _ := errgroup.WithContext(context.TODO())
	for i := range prs {
		i := i
		group.Go(func() error {
			prRef, err := listPullRequest(client, tm, prs[i])
			if err != nil {
				return err
			}

			pullRequests[i] = *prRef

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return pullRequests, nil
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
		Title:     pr.Title,
		Link:      pr.Links.Self[0].Href,
		Order:     pr.UpdatedDate,
		Reviewers: pr.Reviewers,
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
		if latestComment.UpdatedAt().Before(tm) {
			continue
		}

		prRef.Activities = append(prRef.Activities, Comment{
			Comment:   commentActivity,
			Highlight: latestComment.Author.Slug != pr.Author.User.Slug,
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
	Reviewers  []api.Reviewer
	Activities []Comment
	Highlight  bool
}

func (pr PullRequest) Fprint(w io.Writer) {
	strfmt.Fprintf(w, strfmt.StyleBold, "%s\n", pr.Title)
	fmt.Fprintf(w, "%s\n", pr.Link)
	pr.participantsOverview().Fprint(w)

	fmt.Fprint(w, "\n\n")

	for i := range pr.Activities {
		pr.Activities[i].Fprint(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
	}
}

type prOverview struct {
	Approved  []string
	NeedsWork []string
}

func (o prOverview) Fprint(w io.Writer) {
	strfmt.Fprintf(w, strfmt.StyleGreen, "Approved: %d - %v\n", len(o.Approved), o.Approved)
	if len(o.NeedsWork) > 0 {
		strfmt.Fprintf(w, strfmt.StyleRed, "Needs work: %d - %v", len(o.NeedsWork), o.NeedsWork)
	}
}

func (pr PullRequest) participantsOverview() (o prOverview) {
	for i := range pr.Reviewers {
		slug := pr.Reviewers[i].User.Slug
		switch pr.Reviewers[i].Status {
		case api.ReviewerStatusApproved:
			o.Approved = append(o.Approved, slug)
		case api.ReviewerStatusNeedsWork:
			o.NeedsWork = append(o.NeedsWork, slug)
		}
	}
	return o
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
