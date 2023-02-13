package posts

import (
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/comment"
	"time"
)

type Vote struct {
	User string `json:"user"`
	Vot  int    `json:"vot"`
}
type Customer struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}
type Post struct {
	Score            int                              `json:"score"`
	Views            int                              `json:"views"`
	Url              string                           `json:"url"`
	Style            string                           `json:"type"`
	Title            string                           `json:"title"`
	Category         string                           `json:"category"`
	Text             string                           `json:"text"`
	Votes            []Vote                           `json:"votes"`
	Author           Customer                         `json:"author"`
	Comments         comment.CommentMemmoryRepositore `json:"comments"`
	Created          time.Time                        `json:"created"`
	UpvotePercentage int                              `json:"upvotePercentage"`
	Id               string                           `json:"id"`
}
type PostRepo interface {
	GetAll() ([]*Post, error)
	Add(NewPost *Post) ([]byte, error)
	UpDownVote(id string, user string, flag int) (bool, error)
	GetCategory(category string) ([]*Post, error)
	Delete(id string) (bool, error)
	GetAllFromUsers(userId string) ([]*Post, error)
	ShowPost(id string) (*Post, error)
}
