// Package repository предоставляет абстракции для хранилища пользователей и функцию для создания нового репозитория.
package repository

import (
	"BFG_auth/repository/etcd"
	"BFG_auth/repository/types"
	"github.com/pkg/errors"
	"time"
)

// UserRepository определяет интерфейс для абстракции механизма хранения пользователей.
type UserRepository interface {
	// Create добавляет новую пару ключ-значение в репозиторий.
	// Принимает:
	// - kv: структура, содержащая ключ и значение
	// Возвращает:
	// - error: ошибка, если таковая имеется
	Create(kv types.KV) error

	// Read извлекает значение по заданному ключу из репозитория.
	// Принимает:
	// - key: ключ для поиска значения
	// Возвращает:
	// - types.KV: структура, содержащая ключ и значение
	// - error: ошибка, если таковая имеется
	Read(key types.Key) (types.KV, error)

	// Update обновляет существующую пару ключ-значение в репозитории.
	// Принимает:
	// - kv: структура, содержащая ключ и значение
	// Возвращает:
	// - error: ошибка, если таковая имеется
	Update(kv types.KV) error

	// Delete удаляет значение по заданному ключу из репозитория.
	// Принимает:
	// - key: ключ для удаления значения
	// Возвращает:
	// - error: ошибка, если таковая имеется
	Delete(key types.Key) error
}

// NewRepository создает и возвращает новый репозиторий на основе заданного имени.
// Принимает:
// - repoName: имя репозитория, который нужно создать (например, "etcd")
// Возвращает:
// - UserRepository: интерфейс для работы с репозиторием
// - error: ошибка, если таковая имеется
func NewRepository(repoName string) (UserRepository, error) {
	switch repoName {
	case "etcd":
		// Создание нового экземпляра репозитория etcd
		return etcd.NewEtcd([]string{"http://localhost:2379"}, 5*time.Second)
	default:
		// Если заданное имя репозитория не поддерживается, возвращается ошибка
		return nil, errors.New("unknown identity_providers type")
	}
}
