package activityconv

import (
	"github.com/auvn/go-atlassian/bitbucket/api/pr/activity"
	"github.com/auvn/go-atlassian/bitbucket/api/pr/prconv"
	"github.com/auvn/go-json/jsonutil"
)

func CommentFromObject(obj jsonutil.Object) activity.Comment {
	const (
		keyComment       = "comment"
		keyCommentAnchor = "commentAnchor"
	)

	val, _ := obj.Value(keyComment)
	commentAnchor, _ := obj.Value(keyCommentAnchor)

	return activity.Comment{
		Comment:       prconv.CommentFromObject(val.Object()),
		CommentAnchor: CommentAnchorFromObject(commentAnchor.Object()),
	}
}

func CommentAnchorFromObject(obj jsonutil.Object) activity.CommentAnchor {
	const (
		keyLine = "line"
		keyPath = "path"
	)
	line, _ := obj.Value(keyLine)
	path, _ := obj.Value(keyPath)

	return activity.CommentAnchor{
		Path: path.String(),
		Line: line.Int(),
	}
}

func CommentsFromObjects(objs ...jsonutil.Object) []activity.Comment {
	if objs == nil {
		return nil
	}
	comments := make([]activity.Comment, len(objs))
	for i := range objs {
		comments[i] = CommentFromObject(objs[i])
	}
	return comments
}
