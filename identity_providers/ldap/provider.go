// Package ldap предоставляет реализацию провайдера идентификации на основе LDAP.
package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

// LDAP представляет настройки для подключения к LDAP серверу.
type LDAP struct {
	Server string // Адрес сервера
	Port   string // Порт сервера
	BaseDN string // Базовое доменное имя (DN) для поиска в LDAP
}

// Authenticate выполняет аутентификацию пользователя в LDAP.
// Принимает:
// - user: имя пользователя
// - password: пароль пользователя
// Возвращает:
// - bool: true, если аутентификация успешна; false в противном случае
// - error: ошибка, если таковая имеется
func (l *LDAP) Authenticate(user, password string) (bool, error) {
	fmt.Println("Authenticating using LDAP...")

	// Подключение к LDAP серверу
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", l.Server, l.Port))
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	// Закрытие соединения после завершения функции
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Failed to close connection: %v", err)
		}
	}()

	// Формирование полного DN для пользователя
	userDN := fmt.Sprintf("uid=%s,ou=people,dc=localhost", user)
	fmt.Println(userDN)
	fmt.Println(password)

	// Привязка (Bind) к серверу с указанным DN и паролем
	err = conn.Bind(userDN, password)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// Если выполнение дошло до этого места, это означает, что операция Bind была успешной,
	// и, следовательно, аутентификация пользователя прошла успешно.
	return true, nil
}
