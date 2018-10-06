package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
	"github.com/auvn/go-atlassian/bitbucket/api/pr"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity/activityconv"
	"github.com/auvn/go-atlassian/bitbucket/api/resource"
	"github.com/auvn/go-atlassian/bitbucket/api/resource/dashboard"
	"github.com/auvn/go-json/jsonutil"
	"github.com/auvn/go-shell/output"
	"github.com/auvn/go-shell/strfmt"
	"gopkg.in/yaml.v2"
)

var (
	options = struct {
		ConfigFile    string
		MaxCommentAge time.Duration
	}{}
)

func fatal(err error) {
	log.Fatal(err)
}

func init() {
	flag.StringVar(&options.ConfigFile, "config", ".config", "configuration file")
	flag.DurationVar(&options.MaxCommentAge, "age", 0, "max comment age")
}

type Config struct {
	Auth string
	URL  string
}

func config() (cfg Config) {
	bb, err := ioutil.ReadFile(options.ConfigFile)
	if err != nil {
		fatal(err)
	}

	if err := yaml.Unmarshal(bb, &cfg); err != nil {
		fatal(err)
	}

	return cfg
}

func main() {
	flag.Parse()

	cfg := config()

	client := &atlassian.RestClient{
		Client: atlassian.Client{
			Auth:    cfg.Auth,
			BaseURL: cfg.URL,
		},
	}

	bitbucketAPI := resource.Latest{}

	prs, err := allPullRequests(client,
		dashboard.New(bitbucketAPI).
			PullRequests().
			WithParams().
			WithOrder().Newest().
			WithRole().Author().
			WithState().Open())
	if err != nil {
		fatal(err)
	}

	var pullRequests pullRequests

	for i := range prs {
		repo := prs[i].FromRef.Repository
		activities, err := allActivities(client,
			bitbucketAPI.
				Project(repo.Project.Key).
				Repo(repo.Slug).
				PullRequest(prs[i].ID).
				Activities())
		if err != nil {
			fatal(err)
		}

		if len(activities) == 0 {
			continue
		}

		prRef := pullRequest{
			Title:         prs[i].Title,
			Link:          prs[i].Links.Self[0].Href,
			Author:        prs[i].Author.User.EmailAddress,
			MaxCommentAge: options.MaxCommentAge,
		}

		for j := len(activities) - 1; j >= 0; j-- {
			if activity.ActionOf(activities[j]) != activity.ActionCommented {
				continue
			}

			prRef.Activities = append(
				prRef.Activities, prCommentActivity{
					Comment: activityconv.CommentFromObject(activities[j]),
				})
		}
		pullRequests.PRs = append(pullRequests.PRs, prRef)

	}

	pullRequests.Fprintf(os.Stdout)
}

func allActivities(client *atlassian.RestClient, resource api.Resource) ([]jsonutil.Object, error) {
	objects := make([]jsonutil.Object, 0)
	path := resource.Path()
	for {
		resp, err := activity.GetPage(context.TODO(), client, path)
		if err != nil {
			return nil, err
		}

		objects = append(objects, resp.Values...)

		if resp.IsLastPage {
			break
		}

		path = resp.Page.Next(path)
	}

	return objects, nil
}

func allPullRequests(client *atlassian.RestClient, resource api.Resource) ([]api.PullRequest, error) {
	prs := make([]api.PullRequest, 0)
	path := resource.Path()
	for {
		resp, err := pr.GetPage(context.TODO(), client, path)
		if err != nil {
			return nil, err
		}

		prs = append(prs, resp.Values...)

		if resp.IsLastPage {
			break
		}

		path = resp.Page.Next(path)
	}
	return prs, nil
}

type pullRequests struct {
	PRs []pullRequest
}

func (prs pullRequests) Fprintf(w io.Writer) {
	for i := range prs.PRs {
		prs.PRs[i].Fprint(w)
		fmt.Fprintln(w)
	}

}

type pullRequest struct {
	Title         string
	Link          string
	Author        string
	Activities    []prCommentActivity
	MaxCommentAge time.Duration
}

func (pr pullRequest) Fprint(w io.Writer) {
	strfmt.Fprintf(w, strfmt.StyleBold, "%s\n", pr.Title)
	fmt.Fprintf(w, "%s\n\n", pr.Link)

	now := time.Now().UTC()

	for i := range pr.Activities {
		prefixWriter := w

		latestComment := pr.Activities[i].LatestComment()

		if pr.MaxCommentAge > 0 &&
			latestComment.UpdatedAt().Add(pr.MaxCommentAge).Before(now) {
			continue
		}

		if latestComment.Author.EmailAddress == pr.Author {
			prefixWriter = output.NewPrefixWriter(w, "        ")
		} else {
			prefixWriter = output.NewPrefixWriter(w, strfmt.Bold("  =/=/= "))
		}

		pr.Activities[i].Fprint(prefixWriter)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
	}
}

type prCommentActivity struct {
	Comment activity.Comment
}

func latestComment(c pr.Comment) (latest pr.Comment) {
	latest = c

	for i := range c.Comments {
		next := latestComment(c.Comments[i])
		if next.UpdatedDate >= latest.UpdatedDate {
			latest = next
		}
	}

	return latest
}

func (a prCommentActivity) LatestComment() pr.Comment {
	return latestComment(a.Comment.Comment)
}

func (a prCommentActivity) Fprint(w io.Writer) {
	if a.Comment.CommentAnchor.Path != "" {
		strfmt.Fprintf(w, strfmt.StyleBold, "Line: %d Path: %s:\n", a.Comment.CommentAnchor.Line, a.Comment.CommentAnchor.Path)
	}
	prComments{Comments: []pr.Comment{a.Comment.Comment}}.Fprint(w)
}

type prComments struct {
	indentWidth int
	Comments    []pr.Comment
}

func (cc prComments) Fprint(w io.Writer) {
	indentWriter := output.NewSameIndentWriter(w, cc.indentWidth)
	for _, c := range cc.Comments {
		fmt.Fprintf(indentWriter, "%s: %s\n", strfmt.Bold(c.Author.EmailAddress), c.Text)

		prComments{indentWidth: cc.indentWidth + 2, Comments: c.Comments}.Fprint(w)
	}
}
