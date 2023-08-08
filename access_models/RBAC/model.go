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
		return errors.New(fmt.Sprintf("No %s in repo", uid)), nil, nil
	}
	return nil, kv.Key, kv.Value
}

func (s ServiceRBAC) writeToRepo(kv repoTypes.KV) error {
	//идем в репо ищем юзера
	err := s.Repo.Create(kv)
	if err != nil {
		return err
	}
	return nil
}

// DefineUser defines a new user.
func (s ServiceRBAC) CreateUser(uid, idp string, newRoles, newGroups []string) (*accessTypes.User, error) {

	user := accessTypes.User{
		Uid: accessTypes.Uid(uid),
		IdP: accessTypes.IdP(idp),
	}
	user.IdRoles = make([]accessTypes.IdRole, len(newRoles))
	user.Roles = make(map[accessTypes.IdRole]*accessTypes.Role, len(newRoles))
	for i, role := range newRoles {
		readRole, err := s.ReadRole(accessTypes.Uid(role))
		if err != nil {
			return nil, err
		}
		user.IdRoles[i] = readRole.Id
		user.Roles[readRole.Id] = readRole
	}

	user.IdGroups = make([]accessTypes.IdGroup, len(newGroups))
	user.Groups = make(map[accessTypes.IdGroup]*accessTypes.Group, len(newGroups))
	for i, group := range newGroups {
		readGroup, err := s.ReadGroup(accessTypes.Uid(group))
		if err != nil {
			return nil, err
		}
		user.IdGroups[i] = readGroup.Id
		user.Groups[readGroup.Id] = readGroup
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonUser})
	if err != nil {
		return nil, err
	}

	users[string(user.Uid)] = &user

	return &user, nil
}

// DefineRole defines a new role.
func (s ServiceRBAC) CreateRole(id string,
	newPermissions map[accessTypes.IdObject]map[accessTypes.PermissionEnum]bool) (*accessTypes.Role, error) {
	///TODO: check newPermissions, if not correct return error

	idObjects := make([]accessTypes.IdObject, len(newPermissions))

	i := 0
	for idObject := range newPermissions {
		idObjects[i] = idObject
		i++
	}

	role := accessTypes.Role{
		Id:        accessTypes.IdRole(id),
		IdObjects: idObjects,
	}

	role.Permissions = make(map[accessTypes.IdObject]map[accessTypes.PermissionEnum]bool, len(newPermissions))
	role.Permissions = newPermissions

	jsonRole, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(id), Value: jsonRole})
	if err != nil {
		return nil, err
	}

	roles[string(role.Id)] = &role

	return &role, nil
}

// DefineGroup defines a new group.
func (s ServiceRBAC) CreateGroup(id string, newRoles []string) (*accessTypes.Group, error) {
	group := accessTypes.Group{
		Id: accessTypes.IdGroup(id),
	}

	group.IdRoles = make([]accessTypes.IdRole, len(newRoles))
	group.Roles = make(map[accessTypes.IdRole]*accessTypes.Role, len(newRoles))
	for i, role := range newRoles {
		readRole, err := s.ReadRole(accessTypes.Uid(role))
		if err != nil {
			return nil, err
		}
		group.IdRoles[i] = readRole.Id
		group.Roles[readRole.Id] = readRole
	}

	jsonGroup, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(id), Value: jsonGroup})
	if err != nil {
		return nil, err
	}

	groups[string(group.Id)] = &group

	return &group, nil
}

// DefineObject defines a new object.
func (s ServiceRBAC) CreateObject(id accessTypes.IdObject) (*accessTypes.Object, error) {
	var object accessTypes.Object
	object.Id = id

	jsonObject, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(id), Value: jsonObject})
	if err != nil {
		return nil, err
	}

	//objects[id] = &object
	return &object, nil
}

func (s ServiceRBAC) ReadUser(uid accessTypes.Uid) (*accessTypes.User, error) {
	var user accessTypes.User

	err, id, attributes := s.readFromRepo(uid)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(attributes, &user)

	if err != nil {
		return nil, err
	}

	for _, role := range user.IdRoles {
		_, err := s.ReadRole(accessTypes.Uid(role))
		if err != nil {
			return nil, err
		}
	}

	for _, group := range user.IdGroups {
		_, err := s.ReadGroup(accessTypes.Uid(group))
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("RBAC: read user %s: %s\n", id, string(attributes))
	user.Uid = accessTypes.Uid(id)
	users[string(id)] = &user

	return &user, nil
}

func (s ServiceRBAC) ReadRole(uid accessTypes.Uid) (*accessTypes.Role, error) {
	err, id, attributes := s.readFromRepo(uid)

	var role accessTypes.Role
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(attributes, &role)

	if err != nil {
		return nil, err
	}

	fmt.Printf("RBAC: read role %s: %s\n", id, string(attributes))
	roles[string(id)] = &role
	return &role, nil
}

func (s ServiceRBAC) ReadGroup(uid accessTypes.Uid) (*accessTypes.Group, error) {
	var group accessTypes.Group

	err, id, attributes := s.readFromRepo(uid)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(attributes, &group)

	if err != nil {
		return nil, err
	}

	fmt.Printf("RBAC: read group %s: %s\n", id, string(attributes))
	groups[string(id)] = &group
	return &group, nil
}

func (s ServiceRBAC) ReadObject(uid accessTypes.Uid) (*accessTypes.Object, error) {
	err, id, attributes := s.readFromRepo(uid)

	var object accessTypes.Object
	if err != nil {
		return nil, err
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
	jsonUser, err := json.Marshal(newUser)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(newUser.Uid), Value: jsonUser})
	if err != nil {
		return false, err
	}
	// Обновляем пользователя в локальной мапе
	users[string(uid)] = newUser
	return true, nil
}

func (s ServiceRBAC) UpdateRole(uid accessTypes.Uid, newRole *accessTypes.Role) (bool, error) {
	// Обновляем роль в репозитории
	jsonRole, err := json.Marshal(newRole)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonRole})
	if err != nil {
		return false, err
	}
	// Обновляем роль в локальной мапе
	roles[string(uid)] = newRole
	return true, nil
}

func (s ServiceRBAC) UpdateGroup(uid accessTypes.Uid, newGroup *accessTypes.Group) (bool, error) {
	// Обновляем группу в репозитории
	jsonGroup, err := json.Marshal(newGroup)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonGroup})
	if err != nil {
		return false, err
	}
	// Обновляем группу в локальной мапе
	groups[string(uid)] = newGroup
	return true, nil
}

func (s ServiceRBAC) UpdateObject(uid accessTypes.Uid, newObject *accessTypes.Object) (bool, error) {
	// Обновляем объект в репозитории
	jsonObject, err := json.Marshal(newObject)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonObject})
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
