package repository

import (
	"BFG_auth/repository/etcd"
	"BFG_auth/repository/types"
	"github.com/pkg/errors"
	"time"
)

// UserRepository abstracts the storage mechanism for Users
type UserRepository interface {
	Create(kv types.KV) error
	Read(key types.Key) (types.KV, error)
	Update(kv types.KV) error
	Delete(key types.Key) error
}

func NewRepository(repoName string) (UserRepository, error) {
	switch repoName {
	case "etcd":
		return etcd.NewEtcd([]string{"http://localhost:2379"}, 5*time.Second)
	default:
		return nil, errors.New("unknown identity_providers type")
	}
}
