package bitbucketutil

import (
	"context"

	"github.com/auvn/go-atlassian/atlassian"
	"github.com/auvn/go-atlassian/bitbucket/api"
	"github.com/auvn/go-atlassian/bitbucket/api/pr"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity"
	"github.com/auvn/go-json/jsonutil"
)

func GetActivities(client *atlassian.DefaultClient, resource api.Resource) ([]jsonutil.Object, error) {
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

func GetPullRequests(client *atlassian.DefaultClient, resource api.Resource) ([]api.PullRequest, error) {
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

func LatestComment(c pr.Comment) (latest pr.Comment) {
	latest = c

	for i := range c.Comments {
		next := LatestComment(c.Comments[i])
		if next.UpdatedDate >= latest.UpdatedDate {
			latest = next
		}
	}

	return latest
}
