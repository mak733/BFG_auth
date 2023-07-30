package token_service

import (
	"BFG_auth/token_service/JWT"
	"fmt"
	"github.com/pkg/errors"
)

type TokenManager interface {
	GenerateToken(username string) (string, error)
	ValidateToken(tokenString string) (string, error)
}

func CreateTokenManager(managerName string) (TokenManager, error) {
	switch managerName {
	case "JWC":
		fmt.Printf("Get %s\n", managerName)
		return JWT.NewTokenService(), nil
	default:
		return nil, errors.New("unknown token manager type")
	}
}
