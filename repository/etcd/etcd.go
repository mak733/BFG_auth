package etcd

import (
	"BFG_auth/model"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
)

// EtcdUserRepository is an etcd-based User repository
type EtcdUserRepository struct {
	Client *clientv3.Client
}

func (repo *EtcdUserRepository) Create(user model.User) error {
	_, err := repo.Client.Put(context.Background(), fmt.Sprintf("%d", user.Id),
		fmt.Sprintf("%s:%d", user.Name, user.RoleIDs))
	return err
}

func (repo *EtcdUserRepository) Read(id model.IdUser) (*model.User, error) {
	resp, err := repo.Client.Get(context.Background(), string(id))
	if err != nil {
		return nil, err
	}
	for _, ev := range resp.Kvs {
		parts := strings.Split(string(ev.Value), ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid data format for user")
		}
		return &model.User{}, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (repo *EtcdUserRepository) Update(user model.User) error {
	_, err := repo.Client.Put(context.Background(), fmt.Sprintf("%d", user.Id),
		fmt.Sprintf("%s"))
	return err
}

func (repo *EtcdUserRepository) Delete(id model.IdUser) error {
	_, err := repo.Client.Delete(context.Background(), string(id))
	return err
}

func NewREtcd() (*EtcdUserRepository, error) {
	return &EtcdUserRepository{Client: nil}, nil
}
