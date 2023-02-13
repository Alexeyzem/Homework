package comment

import "gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/posts"

type CommentMemmoryRepositore struct {
	data   []*Comment
	lastId string
}

func (receiver *CommentMemmoryRepositore) NewMemmmoryRepo() *CommentMemmoryRepositore {
	return &CommentMemmoryRepositore{
		data: make([]*Comment, 0, 10),
	}
}
func (receiver *CommentMemmoryRepositore) Add(postId string, newComment *Comment) (bool, error) {
	repository := posts.PostMemoryRepository{}
	for _, val := range repository.Data {
		if val.Id == postId {
			val.Comments.data = append(val.Comments.data, newComment)
			return true, nil
		}
	}
	return false, nil
}
func (receiver *CommentMemmoryRepositore) Delete(postId, commentId string) (bool, error) {
	repository := posts.PostMemoryRepository{}
	for _, val := range repository.Data {
		if val.Id == postId {
			for i, tmp := range val.Comments.data {
				if tmp.Id == commentId {
					val.Comments.data = append(val.Comments.data[:i], val.Comments.data[i+1:]...)
					return true, nil
				}
			}
		}
	}
	return false, nil
}
