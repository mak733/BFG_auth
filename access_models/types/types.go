package types

type PermissionEnum int
type IdUser int64
type IdRole int64
type IdGroup int64
type IdObject int64

const (
	Read PermissionEnum = iota
	Write
	Execute
)

func (p PermissionEnum) String() string {
	return [...]string{"Read", "Write", "Execute"}[p]
}
