package apiconv

import (
	"github.com/auvn/go-atlassian/bitbucket/api"
	"github.com/auvn/go-json/jsonutil"
)

func UserFromObject(obj jsonutil.Object) api.User {
	const (
		keyEmailAddress = "emailAddress"
		keySlug         = "slug"
	)
	email, _ := obj.Value(keyEmailAddress)
	slug, _ := obj.Value(keySlug)
	return api.User{
		EmailAddress: email.String(),
		Slug:         slug.String(),
	}
}
