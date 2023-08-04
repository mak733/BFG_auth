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
