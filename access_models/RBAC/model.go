package RBAC

import (
	"BFG_auth/access_models/accessTypes"
	"BFG_auth/repository"
	repoTypes "BFG_auth/repository/types"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type Group struct {
	Id      accessTypes.IdGroup
	Name    string
	RoleIDs []accessTypes.IdRole
}

type Role struct {
	Id          accessTypes.IdRole
	Name        string
	Permissions []accessTypes.PermissionEnum // Список ID разрешений, связанных с ролью
}

type Object struct {
	Id   accessTypes.IdObject
	Name string
}

var (
	users   = make(map[string]*accessTypes.User)
	groups  = make(map[string]*Group)
	roles   = make(map[string]*Role)
	objects = make(map[string]*Object)
)

// RBACService represents the Role-Based Access Control ServiceRBAC
type ServiceRBAC struct {
	Repo repository.UserRepository
}

func (s ServiceRBAC) readFromRepo(uid accessTypes.Uid) (error, repoTypes.Key, repoTypes.Value) {
	//идем в репо ищем юзера
	kv, err := s.Repo.Read(repoTypes.Key(uid))
	if err != nil {
		return errors.New(fmt.Sprintf("No username %s in repo", uid)), nil, nil
	}
	return nil, kv.Key, kv.Value
}

// DefineUser defines a new user.
func (s ServiceRBAC) CreateUser(username string) (*accessTypes.User, error) {
	return nil, nil
}

// DefineRole defines a new role.
func (s ServiceRBAC) CreateRole(id accessTypes.IdRole, permissions []byte) error {
	var role Role

	json.Unmarshal(permissions, &role)
	role.Id = id
	//roles[id] = &role

	return nil
}

// DefineGroup defines a new group.
func (s ServiceRBAC) CreateGroup(id accessTypes.IdGroup, permissions []byte) error {
	var group Group
	json.Unmarshal(permissions, &group)
	group.Id = id
	//groups[id] = &group
	return nil
}

// DefineObject defines a new object.
func (s ServiceRBAC) CreateObject(id accessTypes.IdObject, attributes []byte) error {
	var object Object
	json.Unmarshal(attributes, &object)
	object.Id = id
	//objects[id] = &object
	return nil
}

func (s ServiceRBAC) ReadUser(uid accessTypes.Uid) (*accessTypes.User, error) {
	err, id, attributes := s.readFromRepo(uid)

	var user accessTypes.User
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("RBAC: create %s: %s\n", id, string(attributes))
	err = json.Unmarshal(attributes, &user)

	if err != nil {
		return nil, err
	}

	user.Uid = accessTypes.Uid(id)
	users[string(id)] = &user
	return users[string(id)], nil
}
