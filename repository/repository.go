package repository

import (
	"BFG_auth/model"
	"BFG_auth/repository/etcd"
	"github.com/pkg/errors"
)

type User model.User

// UserRepository abstracts the storage mechanism for Users
type UserRepository interface {
	Create(user model.User) error
	Read(id model.IdUser) (*model.User, error)
	Update(user model.User) error
	Delete(id model.IdUser) error
}

func NewRepository(repoName string) (UserRepository, error) {
	switch repoName {
	case "etcd":
		return etcd.NewREtcd()
	default:
		return nil, errors.New("unknown IdP type")
	}
}
