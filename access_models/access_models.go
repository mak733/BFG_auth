package access_models

import (
	"BFG_auth/access_models/RBAC"
	"BFG_auth/access_models/accessTypes"
	"BFG_auth/repository"
	"fmt"
	"github.com/pkg/errors"
)

type AccessControl interface {
	CreateUser(username string) (*accessTypes.User, error)
	CreateRole(id accessTypes.IdRole, permissions []byte) error
	CreateGroup(id accessTypes.IdGroup, permissions []byte) error
	CreateObject(id accessTypes.IdObject, attributes []byte) error

	ReadUser(id accessTypes.Uid) (*accessTypes.User, error)
	ReadRole(id accessTypes.Uid) (*accessTypes.Role, error)
	ReadGroup(id accessTypes.Uid) (*accessTypes.Group, error)
	ReadObject(id accessTypes.Uid) (*accessTypes.Object, error)

	UpdateUser(id accessTypes.Uid, newUser *accessTypes.User) (bool, error)
	UpdateRole(id accessTypes.Uid, newRole *accessTypes.Role) (bool, error)
	UpdateGroup(id accessTypes.Uid, newGroup *accessTypes.Group) (bool, error)
	UpdateObject(id accessTypes.Uid, newObject *accessTypes.Object) (bool, error)

	DeleteUser(id accessTypes.Uid) (bool, error)
	DeleteRole(id accessTypes.Uid) (bool, error)
	DeleteGroup(id accessTypes.Uid) (bool, error)
	DeleteObject(id accessTypes.Uid) (bool, error)
}

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
