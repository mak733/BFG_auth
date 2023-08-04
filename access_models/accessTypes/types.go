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

	Roles  IdRole  `json:"Role"`
	Groups IdGroup `json:"Group"`
}
