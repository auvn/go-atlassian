package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api/pr"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity/activityconv"
	"github.com/auvn/go-atlassian/bitbucket/api/resource"
	"github.com/auvn/go-atlassian/bitbucket/api/resource/dashboard"
	"github.com/auvn/go-shell/output"
	"github.com/auvn/go-shell/strfmt"
	"gopkg.in/yaml.v2"
)

var (
	ctx     = context.TODO()
	options = struct {
		ConfigFile string
	}{}
)

func fatal(err error) {
	log.Fatal(err)
}

func init() {
	flag.StringVar(&options.ConfigFile, "config", ".config", "configuration file")
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

	resp, err := pr.GetPage(ctx, client,
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

	for i := range resp.Values {
		repo := resp.Values[i].FromRef.Repository
		aResp, err := activity.Get(ctx, client,
			bitbucketAPI.
				Project(repo.Project.Key).
				Repo(repo.Slug).
				PullRequest(resp.Values[i].ID).
				Activities())
		if err != nil {
			fatal(err)
		}

		if len(aResp.Values) == 0 {
			continue
		}

		prRef := pullRequest{
			Title:  resp.Values[i].Title,
			Link:   resp.Values[i].Links.Self[0].Href,
			Author: resp.Values[i].Author.User.EmailAddress,
		}

		for j := len(aResp.Values) - 1; j >= 0; j-- {
			if activity.ActionOf(aResp.Values[j]) != activity.ActionCommented {
				continue
			}

			prRef.Activities = append(
				prRef.Activities, prCommentActivity{
					Comment: activityconv.CommentFromObject(aResp.Values[j]),
				})
		}
		pullRequests.PRs = append(pullRequests.PRs, prRef)

	}

	pullRequests.Fprintf(os.Stdout)
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
	Title      string
	Link       string
	Author     string
	Activities []prCommentActivity
}

func (pr pullRequest) Fprint(w io.Writer) {
	strfmt.Fprintf(w, strfmt.StyleBold, "%s\n", pr.Title)
	fmt.Fprintf(w, "%s\n\n", pr.Link)

	for i := range pr.Activities {
		prefixWriter := w

		if pr.Activities[i].IsLatestCommentBy(pr.Author) {
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

func walkComment(level int, c pr.Comment, fn func(s string) bool) (maxLevel int, ok bool) {
	for i := range c.Comments {
		maxLevel, ok = walkComment(level+1, c.Comments[i], fn)
		if ok {
			return maxLevel, true
		}
	}

	if fn(c.Author.EmailAddress) && level >= maxLevel {
		return level, true
	}
	return maxLevel, false
}

func (a prCommentActivity) IsLatestCommentBy(emailAddress string) bool {
	_, ok := walkComment(0, a.Comment.Comment, func(str string) bool {
		return emailAddress == str
	})

	return ok
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
