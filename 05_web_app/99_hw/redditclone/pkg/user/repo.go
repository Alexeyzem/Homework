package user

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoUser  = errors.New("No user found")
	ErrBadPass = errors.New("Invald password")
)

type UserMemoryRepository struct {
	data   map[string]*User
	lastId uint32
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: map[string]*User{
			"alexeyZ admin": {
				ID:       1,
				Login:    "alexeyZ admin",
				password: "idk what i am doing. WTF",
			},
		},
		lastId: 1,
	}
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	if u.password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}
func (repo *UserMemoryRepository) Registrate(login, pass string) (bool, error) {
	if strings.EqualFold(login, pass) {
		return false, fmt.Errorf("password and login must match")
	}
	repo.lastId++
	repo.data[login] = &User{ID: repo.lastId, Login: login, password: pass}
	return true, nil
}
