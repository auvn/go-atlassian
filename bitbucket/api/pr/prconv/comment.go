package prconv

import (
	"github.com/auvn/go-atlassian/bitbucket/api/apiconv"
	"github.com/auvn/go-atlassian/bitbucket/api/pr"
	"github.com/auvn/go-json/jsonutil"
)

func CommentFromObject(obj jsonutil.Object) pr.Comment {
	const (
		keyText        = "text"
		keyAuthor      = "author"
		keyComments    = "comments"
		keyUpdatedDate = "updatedDate"
	)

	text, _ := obj.Value(keyText)
	author, _ := obj.Value(keyAuthor)
	comments, _ := obj.Value(keyComments)
	updatedDate, _ := obj.Value(keyUpdatedDate)

	return pr.Comment{
		Author:      apiconv.UserFromObject(author.Object()),
		Text:        text.String(),
		Comments:    CommentsFromObjects(comments.Objects()...),
		UpdatedDate: updatedDate.Int64(),
	}
}

func CommentsFromObjects(objs ...jsonutil.Object) []pr.Comment {
	if objs == nil {
		return nil
	}

	comments := make([]pr.Comment, len(objs))
	for i := range objs {
		comments[i] = CommentFromObject(objs[i])
	}
	return comments
}
