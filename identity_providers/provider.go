package identity_providers

import (
	"BFG_auth/identity_providers/ldap"
	"errors"
)

// IdP - Interface for Identity Providers
type IdP interface {
	Authenticate(user, password string) (bool, error)
}

func NewIdp(idpType string) (IdP, error) {
	switch idpType {
	case "ldap":
		return &ldap.LDAP{
			Server: "localhost",
			Port:   "389",
			BaseDN: "cn=admin,dc=localhost",
		}, nil
	default:
		return nil, errors.New("unknown identity_providers type")
	}
}
