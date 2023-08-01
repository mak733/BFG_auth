package JWT

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("secret")

type Claims struct {
	Username string
	jwt.StandardClaims
}

type TokenService struct {
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (s TokenService) ValidateToken(username, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return false, errors.New("invalid token")
	}

	// проверяем, что имя пользователя в токене совпадает с именем пользователя, для которого мы генерировали токен
	if claims.Username != username {
		return false, errors.New("invalid token: username mismatch")
	}

	return true, nil
}
