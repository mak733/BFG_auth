package RBAC

import (
	"BFG_auth/access_models/accessTypes"
	"BFG_auth/repository"
	repoTypes "BFG_auth/repository/types"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

var (
	users   = make(map[string]*accessTypes.User)
	roles   = make(map[string]*accessTypes.Role)
	groups  = make(map[string]*accessTypes.Group)
	objects = make(map[string]*accessTypes.Object)
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
	var role accessTypes.Role

	json.Unmarshal(permissions, &role)
	role.Id = id
	//roles[id] = &role

	return nil
}

// DefineGroup defines a new group.
func (s ServiceRBAC) CreateGroup(id accessTypes.IdGroup, permissions []byte) error {
	var group accessTypes.Group
	json.Unmarshal(permissions, &group)
	group.Id = id
	//groups[id] = &group
	return nil
}

// DefineObject defines a new object.
func (s ServiceRBAC) CreateObject(id accessTypes.IdObject, attributes []byte) error {
	var object accessTypes.Object
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

	fmt.Printf("RBAC: read %s: %s\n", id, string(attributes))
	err = json.Unmarshal(attributes, &user)

	if err != nil {
		return nil, err
	}

	for _, role := range user.Roles {
		s.ReadRole(accessTypes.Uid(role))
	}

	for _, group := range user.Groups {
		s.ReadGroup(accessTypes.Uid(group))
	}

	for _, object := range user.Objects {
		s.ReadGroup(accessTypes.Uid(object))
	}

	user.Uid = accessTypes.Uid(id)
	users[string(id)] = &user

	return &user, nil
}

func (s ServiceRBAC) ReadRole(uid accessTypes.Uid) (*accessTypes.Role, error) {
	err, id, attributes := s.readFromRepo(uid)

	var role accessTypes.Role
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("RBAC: read %s: %s\n", id, string(attributes))
	err = json.Unmarshal(attributes, &role)

	if err != nil {
		return nil, err
	}

	roles[string(id)] = &role
	return &role, nil
}

func (s ServiceRBAC) ReadGroup(uid accessTypes.Uid) (*accessTypes.Group, error) {
	err, id, attributes := s.readFromRepo(uid)

	var group accessTypes.Group
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("RBAC: read %s: %s\n", id, string(attributes))
	err = json.Unmarshal(attributes, &group)

	if err != nil {
		return nil, err
	}

	groups[string(id)] = &group
	return &group, nil
}

func (s ServiceRBAC) ReadObject(uid accessTypes.Uid) (*accessTypes.Object, error) {
	err, id, attributes := s.readFromRepo(uid)

	var object accessTypes.Object
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("RBAC: read %s: %s\n", id, string(attributes))
	err = json.Unmarshal(attributes, &object)

	if err != nil {
		return nil, err
	}

	objects[string(id)] = &object
	return &object, nil
}

func (s ServiceRBAC) UpdateUser(uid accessTypes.Uid, newUser *accessTypes.User) (bool, error) {
	// Обновляем пользователя в репозитории
	json, err := json.Marshal(newUser)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: json})
	if err != nil {
		return false, err
	}
	// Обновляем пользователя в локальной мапе
	users[string(uid)] = newUser
	return true, nil
}

func (s ServiceRBAC) UpdateRole(uid accessTypes.Uid, newRole *accessTypes.Role) (bool, error) {
	// Обновляем роль в репозитории
	json, err := json.Marshal(newRole)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: json})
	if err != nil {
		return false, err
	}
	// Обновляем роль в локальной мапе
	roles[string(uid)] = newRole
	return true, nil
}

func (s ServiceRBAC) UpdateGroup(uid accessTypes.Uid, newGroup *accessTypes.Group) (bool, error) {
	// Обновляем группу в репозитории
	json, err := json.Marshal(newGroup)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: json})
	if err != nil {
		return false, err
	}
	// Обновляем группу в локальной мапе
	groups[string(uid)] = newGroup
	return true, nil
}

func (s ServiceRBAC) UpdateObject(uid accessTypes.Uid, newObject *accessTypes.Object) (bool, error) {
	// Обновляем объект в репозитории
	json, err := json.Marshal(newObject)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: json})
	if err != nil {
		return false, err
	}
	// Обновляем объект в локальной мапе
	objects[string(uid)] = newObject
	return true, nil
}

func (s ServiceRBAC) DeleteUser(uid accessTypes.Uid) (bool, error) {
	// Удаляем пользователя из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}
	// Удаляем пользователя из локальной мапы
	delete(users, string(uid))
	return true, nil
}

func (s ServiceRBAC) DeleteRole(uid accessTypes.Uid) (bool, error) {
	// Удаляем роль из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}
	// Удаляем роль из локальной мапы
	delete(roles, string(uid))
	return true, nil
}

func (s ServiceRBAC) DeleteGroup(uid accessTypes.Uid) (bool, error) {
	// Удаляем группу из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}
	// Удаляем группу из локальной мапы
	delete(groups, string(uid))
	return true, nil
}

func (s ServiceRBAC) DeleteObject(uid accessTypes.Uid) (bool, error) {
	// Удаляем объект из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}
	// Удаляем объект из локальной мапы
	delete(objects, string(uid))
	return true, nil
}
