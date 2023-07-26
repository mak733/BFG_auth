package model

import (
	"BFG_auth/IdP"
	"time"
)

var (
	dialTimeout = 5 * time.Second
	endPoint    = []string{"localhost:2379"}
)

type IdUser int64
type IdRole int64
type IdGroup int64
type IdObject int64

type PermissionEnum int

const (
	Read PermissionEnum = iota
	Write
	Execute
)

func (p PermissionEnum) String() string {
	return [...]string{"Read", "Write", "Execute"}[p]
}

type User struct {
	Id       IdUser
	Name     string
	RoleIDs  []IdRole  // Список ID ролей, связанных с пользователем
	GroupIDs []IdGroup // Список ID групп, связанных с пользователем
	IdP      IdP.IdP
}

type Group struct {
	Id      IdGroup
	Name    string
	RoleIDs []IdRole
}

type Role struct {
	Id          IdRole
	Name        string
	Permissions []PermissionEnum // Список ID разрешений, связанных с ролью
}

type Object struct {
	Id   IdObject
	Name string
	//TODO: some objects
}
