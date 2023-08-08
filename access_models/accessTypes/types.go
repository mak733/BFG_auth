package accessTypes

import "fmt"

type Uid string
type IdP string

type PermissionEnum string

type Objects []string

type IdRole string
type IdGroup string
type IdObject string

type User struct {
	Uid      Uid       `json:"Uid"`
	IdP      IdP       `json:"IdP"`
	IdRoles  []IdRole  `json:"Role"`
	IdGroups []IdGroup `json:"Group"`

	Roles  map[IdRole]*Role
	Groups map[IdGroup]*Group
}

type Role struct {
	Id        IdRole     `json:"Id"`
	IdObjects []IdObject `json:"IdObjects"`

	Permissions map[IdObject]map[PermissionEnum]bool
}

type Group struct {
	Id      IdGroup  `json:"Id"`
	IdRoles []IdRole `json:"Role"`

	Roles map[IdRole]*Role
}

type Object struct {
	Id IdObject `json:"Id"`
}

func (u *User) CheckPermission(object IdObject, permission PermissionEnum) (bool, error) {

	for _, role := range u.Roles {
		ok := role.CheckPermission(object, permission)
		if ok == true {
			return ok, nil
		}
	}

	for _, group := range u.Groups {
		ok, err := group.CheckPermission(object, permission)
		if err != nil {
			return false, err
		}
		if ok == true {
			return ok, nil
		}
	}

	return false, nil
}

func (r *Role) CheckPermission(object IdObject, permission PermissionEnum) bool {
	if perms, exists := r.Permissions[object]; exists {
		if _, hasPermission := perms[permission]; hasPermission {
			return true
		}
	} else {
		return false
	}
	return false
}

func (g *Group) CheckPermission(object IdObject, permission PermissionEnum) (bool, error) {
	for _, role := range g.Roles {
		ok := role.CheckPermission(object, permission)
		if ok == true {
			return ok, nil
		}
	}

	return false, fmt.Errorf("permission %s not found for object %s", permission, object)
}
