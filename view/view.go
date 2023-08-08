// Package view предоставляет интерфейсы и функции для создания и управления представлениями.
package view

import (
	"BFG_auth/session"
	"BFG_auth/view/http"
	"github.com/pkg/errors"
)

// View определяет интерфейс для представлений, которые могут быть использованы для взаимодействия с пользователем.
type View interface {
	// StartServer запускает сервер представления на указанном адресе с использованием переданного менеджера сессий.
	StartServer(address string, sessionManager *session.Manager) error
}

// NewView создает и возвращает новое представление на основе заданного имени.
// Принимает:
// - ViewName: имя представления (например, "http")
// Возвращает:
// - View: интерфейс для работы с представлением
// - error: ошибка, если задан неправильный тип представления
func NewView(ViewName string) (View, error) {
	switch ViewName {
	case "http":
		return &http.ViewHttp{}, nil
	default:
		return nil, errors.New("unknown controllers type")
	}
}
