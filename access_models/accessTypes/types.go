// Package accessTypes определяет основные типы данных и структуры, используемые для контроля доступа.
package accessTypes

import "fmt"

// Uid представляет уникальный идентификатор пользователя.
type Uid string

// IdP представляет идентификатор поставщика (Provider).
type IdP string

// PermissionEnum представляет тип разрешения.
type PermissionEnum string

// Objects представляет набор объектов.
type Objects []string

// IdRole представляет уникальный идентификатор роли.
type IdRole string

// IdGroup представляет уникальный идентификатор группы.
type IdGroup string

// IdObject представляет уникальный идентификатор объекта.
type IdObject string

// User представляет пользователя с его ролями и группами.
type User struct {
	Uid      Uid       `json:"Uid"`
	IdP      IdP       `json:"IdP"`
	IdRoles  []IdRole  `json:"Role"`
	IdGroups []IdGroup `json:"Group"`

	Roles  map[IdRole]*Role
	Groups map[IdGroup]*Group
}

// Role представляет роль с привязанными к ней объектами и разрешениями.
type Role struct {
	Id        IdRole     `json:"Id"`
	IdObjects []IdObject `json:"IdObjects"`

	Permissions map[IdObject]map[PermissionEnum]bool
}

// Group представляет группу пользователей с привязанными к ней ролями.
type Group struct {
	Id      IdGroup  `json:"Id"`
	IdRoles []IdRole `json:"Role"`

	Roles map[IdRole]*Role
}

// Object представляет объект, к которому можно контролировать доступ.
type Object struct {
	Id IdObject `json:"Id"`
}

// CheckPermission проверяет, имеет ли пользователь разрешение на выполнение действия над объектом.
//
// Принимает:
//   - object: объект, к которому применяется разрешение
//   - permission: тип разрешения
//
// Возвращает:
//   - bool: true, если разрешение есть, иначе false
//   - error: ошибка, если таковая имеется
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

// Принимает:
//   - object: объект, к которому применяется разрешение
//   - permission: тип разрешения
//
// Возвращает:
//   - bool: true, если разрешение есть, иначе false
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

// CheckPermission проверяет, имеет ли группа разрешение на выполнение действия над объектом.
//
// Принимает:
//   - object: объект, к которому применяется разрешение
//   - permission: тип разрешения
//
// Возвращает:
//   - bool: true, если разрешение есть, иначе false
//   - error: ошибка, если таковая имеется
func (g *Group) CheckPermission(object IdObject, permission PermissionEnum) (bool, error) {
	for _, role := range g.Roles {
		ok := role.CheckPermission(object, permission)
		if ok == true {
			return ok, nil
		}
	}

	return false, fmt.Errorf("permission %s not found for object %s", permission, object)
}
