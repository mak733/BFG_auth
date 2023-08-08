// Package access_models предоставляет интерфейсы и функции для управления моделями контроля доступа.
package access_models

import (
	"BFG_auth/access_models/RBAC"
	"BFG_auth/access_models/accessTypes"
	"BFG_auth/repository"
	"fmt"
	"github.com/pkg/errors"
)

// AccessControl определяет интерфейс для управления пользователями, ролями, группами и объектами в контексте контроля доступа.
type AccessControl interface {
	// CRUD операции для пользователей
	CreateUser(uid, idp string, newRoles, newGroups []string) (*accessTypes.User, error)
	ReadUser(id accessTypes.Uid) (*accessTypes.User, error)
	UpdateUser(id accessTypes.Uid, newUser *accessTypes.User) (bool, error)
	DeleteUser(id accessTypes.Uid) (bool, error)

	// CRUD операции для ролей
	CreateRole(id string, newPermissions map[accessTypes.IdObject]map[accessTypes.PermissionEnum]bool) (*accessTypes.Role, error)
	ReadRole(id accessTypes.Uid) (*accessTypes.Role, error)
	UpdateRole(id accessTypes.Uid, newRole *accessTypes.Role) (bool, error)
	DeleteRole(id accessTypes.Uid) (bool, error)

	// CRUD операции для групп
	CreateGroup(id string, newRoles []string) (*accessTypes.Group, error)
	ReadGroup(id accessTypes.Uid) (*accessTypes.Group, error)
	UpdateGroup(id accessTypes.Uid, newGroup *accessTypes.Group) (bool, error)
	DeleteGroup(id accessTypes.Uid) (bool, error)

	// CRUD операции для объектов
	CreateObject(id accessTypes.IdObject) (*accessTypes.Object, error)
	ReadObject(id accessTypes.Uid) (*accessTypes.Object, error)
	UpdateObject(id accessTypes.Uid, newObject *accessTypes.Object) (bool, error)
	DeleteObject(id accessTypes.Uid) (bool, error)
}

// NewAccessControlModel создает и возвращает новую модель контроля доступа на основе заданных имени модели и репозитория.
//
// Принимает:
//   - modelName: имя модели контроля доступа (например, "RBAC")
//   - repoName: имя репозитория
//
// Возвращает:
//   - AccessControl: интерфейс для работы с моделью контроля доступа
//   - error: ошибка, если таковая имеется
func NewAccessControlModel(modelName, repoName string) (AccessControl, error) {
	repo, err := repository.NewRepository(repoName)
	if err != nil {
		fmt.Printf("Error make repository %s\n", repoName)
		return nil, err
	}

	switch modelName {
	case "RBAC":
		fmt.Printf("Get %s\n", modelName)
		return &RBAC.ServiceRBAC{Repo: repo}, nil
	default:
		return nil, errors.New("unknown model type")
	}
}
