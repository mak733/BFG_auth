package etcd

import (
	"BFG_auth/repository/types"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
	"time"
)

// UserRepository EtcdUserRepository is an etcd-based User repository
type UserRepository struct {
	Client *clientv3.Client
}

/*
	func (repo *EtcdUserRepository) Create(user model.User) error {
		_, err := repo.Client.Put(context.Background(), fmt.Sprintf("%d", user.Id),
			fmt.Sprintf("%s:%d", user.Name, user.RoleIDs))
		return err
	}
*/
func (repo *UserRepository) Read(id types.Key) (types.KV, error) {

	resp, err := repo.Client.Get(context.Background(), string(id))
	if err != nil {
		return types.KV{}, err
	}

	if resp.Count == 0 {
		return types.KV{}, errors.New("No key in etcd")
	}

	body := resp.Kvs[0]
	fmt.Printf("Get value by key %s: %s\n", body.Key, body.Value)
	return types.KV{Key: body.Key, Value: body.Value}, nil
}

/*
func (repo *EtcdUserRepository) Update(user model.User) error {
	_, err := repo.Client.Put(context.Background(), fmt.Sprintf("%d", user.Id),
		fmt.Sprintf("%s"))
	return err
}

func (repo *EtcdUserRepository) Delete(id model.IdUser) error {
	_, err := repo.Client.Delete(context.Background(), string(id))
	return err
}
*/

func NewEtcd(endpoints []string, dialTimeout time.Duration) (*UserRepository, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	}

	Client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &UserRepository{Client: Client}, nil
}
