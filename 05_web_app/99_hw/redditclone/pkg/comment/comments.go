package comment

import (
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/posts"
	"time"
)

type Comment struct {
	Created time.Time      `json:"created"`
	Author  posts.Customer `json:"author"`
	Body    string         `json:"body"`
	Id      string         `json:"id"`
}
type CommentRepo interface {
	Add(postId string) (bool, error)
	Delete(postId, commentId string) (bool, error)
}
