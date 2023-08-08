// Package types содержит типы данных, используемые в различных частях приложения.
package types

// Credentials представляет собой учетные данные пользователя, включая имя пользователя и пароль.
type Credentials struct {
	Username string // Username представляет собой имя пользователя.
	Password string // Password содержит пароль пользователя.
}
