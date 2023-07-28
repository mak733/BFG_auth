package types

type PermissionEnum int
type Uid string
type IdRole int64
type IdGroup int64
type IdObject int64

type IdP string

const (
	Read PermissionEnum = iota
	Write
	Execute
)

func (p PermissionEnum) String() string {
	return [...]string{"Read", "Write", "Execute"}[p]
}

type Objects []string
type JWTToken string

type User struct {
	Uid         Uid
	Permissions PermissionEnum
	Objects     Objects
	JWTToken    JWTToken
	IdP         IdP
}
