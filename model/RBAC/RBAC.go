package RBAC

import (
	"BFG_auth/model"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type User model.User

// EtcdUserCRUD - etcd implementation of UserRepository
type EtcdUserCRUD struct {
	Client *clientv3.Client
}

// Create user
func (crud *EtcdUserCRUD) Create(user User) error {
	_, err := crud.Client.Put(context.Background(), fmt.Sprintf("%d", user.Id), fmt.Sprintf("%s:%v", user.Name, user.RoleIDs))
	return err
}

// Read user
func (crud *EtcdUserCRUD) Read(id string) (*User, error) {
	resp, err := crud.Client.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}
	for _, ev := range resp.Kvs {
		// Get name and role from value, in real life consider using json or similar
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	return &User{}, nil
}

// Update user
func (crud *EtcdUserCRUD) Update(user User) error {
	_, err := crud.Client.Put(context.Background(), fmt.Sprintf("%d", user.Id), fmt.Sprintf("%s:%v", user.Name, user.RoleIDs))
	return err
}

// Delete user
func (crud *EtcdUserCRUD) Delete(id string) error {
	_, err := crud.Client.Delete(context.Background(), id)
	return err
}
