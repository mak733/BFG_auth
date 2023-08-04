package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

// LDAP - LDAP Identity Provider
type LDAP struct {
	Server string
	Port   string
	BaseDN string
}

func (l *LDAP) Authenticate(user, password string) (bool, error) {
	fmt.Println("Authenticating using LDAP...")

	// Connect to the LDAP server
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", l.Server, l.Port))
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Failed to close connection: %v", err)
		}
	}()

	// Construct the full DN
	userDN := fmt.Sprintf("uid=%s,ou=people,dc=localhost", user)
	fmt.Println(userDN)
	fmt.Println(password)
	// Bind to the server with given DN and password
	err = conn.Bind(userDN, password)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// If we reach here, it means the Bind operation was successful,
	// and hence the user has been authenticated.
	return true, nil
}
