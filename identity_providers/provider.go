// Package identity_providers предоставляет функциональность для работы с различными провайдерами идентификации.
package identity_providers

import (
	"BFG_auth/identity_providers/ldap"
	"errors"
)

// IdP представляет интерфейс для всех поставщиков идентификации.
type IdP interface {
	// Authenticate выполняет аутентификацию пользователя.
	// Принимает:
	// - user: имя пользователя
	// - password: пароль пользователя
	// Возвращает:
	// - bool: true, если аутентификация успешна; false в противном случае
	// - error: ошибка, если таковая имеется
	Authenticate(user, password string) (bool, error)
}

// NewIdp создает и возвращает объект, реализующий интерфейс IdP, на основе указанного типа поставщика идентификации.
// Принимает:
// - idpType: строка, определяющая тип поставщика идентификации
// Возвращает:
// - IdP: объект, реализующий интерфейс IdP
// - error: ошибка, если указанный тип поставщика не поддерживается
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
