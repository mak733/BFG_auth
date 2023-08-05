package accessTypes

type Uid string
type PermissionEnum []string
type Objects []string
type IdP string

type IdRole []string
type IdGroup []string
type IdObject []string

type User struct {
	Uid         Uid            `json:"Uid"`
	Permissions PermissionEnum `json:"Permissions"`
	Objects     Objects        `json:"Objects"`
	IdP         IdP            `json:"IdP"`
	Roles       IdRole         `json:"Role"`
	Groups      IdGroup        `json:"Group"`
}

type Role struct {
	Id          IdRole
	Name        string
	Permissions []PermissionEnum // Список ID разрешений, связанных с ролью
}

type Group struct {
	Id      IdGroup
	Name    string
	RoleIDs []IdRole
}

type Object struct {
	Id   IdObject
	Name string
}
