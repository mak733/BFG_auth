package IdP

import (
	"BFG_auth/IdP/ldap"
	"BFG_auth/IdP/oauth"
	"BFG_auth/IdP/saml"
	"errors"
)

// IdP - Interface for Identity Providers
type IdP interface {
	Authenticate(username, password string) bool
}

// TODO: add configurations for IdPs
func NewIdp(idpType string) (IdP, error) {
	switch idpType {
	case "ldap":
		return &ldap.LDAP{
			Server: "localhost",
			Port:   "389",
			BaseDN: "dc=example,dc=com",
		}, nil
	case "saml":
		return &saml.SAML{
			EntityID: "entityID",
			// ... other fields ...
		}, nil
	case "oauth":
		return &oauth.OAuth{
			ClientID: "clientID",
			// ... other fields ...
		}, nil
	default:
		return nil, errors.New("unknown IdP type")
	}
}
