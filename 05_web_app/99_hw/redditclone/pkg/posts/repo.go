package posts

import (
	"crypto/rand"
	"errors"
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/comment"
)

var errNotFoundPost = errors.New("invalid post id")

type PostMemoryRepository struct {
	lastId []byte
	Data   []*Post
}

func NewMemoryRepo() *PostMemoryRepository {
	return &PostMemoryRepository{
		Data: make([]*Post, 0, 10),
	}
}
func (repo *PostMemoryRepository) GetAll() ([]*Post, error) {
	return repo.Data, nil
}
func (repo *PostMemoryRepository) Add(NewPost *Post) ([]byte, error) {
	//доработать для разных постов
	//с урлом и просто текстом

	repositore := comment.CommentMemmoryRepositore{}
	NewPost.Comments = *repositore.NewMemmmoryRepo()
	repo.lastId = make([]byte, 16)
	_, err := rand.Read(repo.lastId)
	if err != nil {
		return nil, err
	}
	NewPost.Id = string(repo.lastId)
	repo.Data = append(repo.Data, NewPost)
	return repo.lastId, nil
}
func (repo *PostMemoryRepository) UpDownVote(id string, user string, flag int) (bool, error) {
	for _, post := range repo.Data {
		if post.Id == id {
			for j, val := range post.Votes {
				if val.User == user && val.Vot == 1*flag {
					post.Score -= flag
					post.Votes = append(post.Votes[:j], post.Votes[j+1:]...)
					return true, nil
				} else if val.User == user {
					post.Score += 2 * flag
					post.Votes[j].Vot = 1 * flag
					return true, nil
				}
			}
			post.Score += flag
			post.Votes = append(post.Votes, Vote{User: user, Vot: 1 * flag})
			return true, nil
		}
	}
	return false, nil
}
func (repo *PostMemoryRepository) GetCategory(category string) ([]*Post, error) {
	var result []*Post
	for _, val := range repo.Data {
		if val.Category == category {
			result = append(result, val)
		}
	}
	return result, nil
}
func (repo *PostMemoryRepository) Delete(id string) (bool, error) {
	for i, val := range repo.Data {
		if val.Id == id {
			repo.Data = append(repo.Data[:i], repo.Data[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}
func (repo *PostMemoryRepository) GetAllFromUsers(userId string) ([]*Post, error) {
	var result []*Post
	for _, val := range repo.Data {
		if val.Author.Id == userId {
			result = append(result, val)
		}
	}
	return result, nil
}
func (repo *PostMemoryRepository) ShowPost(id string) (*Post, error) {
	for _, val := range repo.Data {
		if val.Id == id {
			val.Views++
			return val, nil
		}
	}
	return nil, errNotFoundPost
}
