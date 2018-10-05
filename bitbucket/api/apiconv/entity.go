package apiconv

import (
	"github.com/auvn/go-atlassian/bitbucket/api"
	"github.com/auvn/go-json/jsonutil"
)

func UserFromObject(obj jsonutil.Object) api.User {
	const (
		keyEmailAddress = "emailAddress"
	)
	email, _ := obj.Value(keyEmailAddress)
	return api.User{
		EmailAddress: email.String(),
	}
}
