// Package JWT предоставляет функциональность для работы с JWT токенами.
package JWT

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// secretKey используется для подписи и проверки JWT токенов.
var secretKey = []byte("secret")

// Claims представляет собой структуру данных, которая хранится в JWT токене.
type Claims struct {
	Username string
	jwt.StandardClaims
}

// TokenService предоставляет методы для генерации и проверки JWT токенов.
type TokenService struct {
}

// NewTokenService создает и возвращает новый экземпляр TokenService.
func NewTokenService() *TokenService {
	return &TokenService{}
}

// GenerateToken генерирует новый JWT токен для указанного имени пользователя.
// Принимает:
// - username: имя пользователя
// Возвращает:
// - string: сгенерированный JWT токен
// - error: ошибка, если таковая имеется
func (s *TokenService) GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken проверяет действительность JWT токена для указанного имени пользователя.
// Принимает:
// - username: имя пользователя
// - tokenString: JWT токен для проверки
// Возвращает:
// - bool: true, если токен действителен, иначе false
// - error: ошибка, если таковая имеется
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
