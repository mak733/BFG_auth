// Package etcd предоставляет реализацию хранилища пользователей на основе etcd.
package etcd

import (
	"BFG_auth/repository/types"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
	"time"
)

// UserRepository - это структура, представляющая репозиторий пользователей на основе etcd.
type UserRepository struct {
	Client *clientv3.Client // Client предоставляет интерфейс для взаимодействия с etcd.
}

// Create создает новую запись в etcd.
//
// Принимает:
//   - kv: пару ключ-значение для сохранения.
//
// Возвращает:
//   - error: ошибка, если таковая имеется.
func (repo *UserRepository) Create(kv types.KV) error {
	_, err := repo.Client.Put(context.Background(), string(kv.Key), string(kv.Value))
	return err
}

// Read читает запись из etcd по заданному ключу.
//
// Принимает:
//   - id: ключ для поиска.
//
// Возвращает:
//   - types.KV: пару ключ-значение.
//   - error: ошибка, если таковая имеется.
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

// Update обновляет запись в etcd.
//
// Принимает:
//   - kv: пару ключ-значение для обновления.
//
// Возвращает:
//   - error: ошибка, если таковая имеется.
func (repo *UserRepository) Update(kv types.KV) error {
	_, err := repo.Client.Put(context.Background(), string(kv.Key), string(kv.Value))
	return err
}

// Delete удаляет запись из etcd по заданному ключу.
//
// Принимает:
//   - key: ключ для удаления.
//
// Возвращает:
//   - error: ошибка, если таковая имеется.
func (repo *UserRepository) Delete(key types.Key) error {
	_, err := repo.Client.Delete(context.Background(), string(key))
	return err
}

// NewEtcd создает новый экземпляр UserRepository с инициализированным etcd клиентом.
//
// Принимает:
//   - endpoints: адреса etcd серверов.
//   - dialTimeout: время ожидания при подключении.
//
// Возвращает:
//   - *UserRepository: указатель на UserRepository.
//   - error: ошибка, если таковая имеется.
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
