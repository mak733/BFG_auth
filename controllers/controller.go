package controllers

import (
	"BFG_auth/controllers/ubuntu"
	"github.com/pkg/errors"
)

// API представляет интерфейс идентификационных провайдеров (IdP).
// Он содержит методы для получения различной информации об идентификационных провайдерах.
type API interface {
	Name() string    // Название провайдера
	Time() string    // Временная метка провайдера
	Disk() string    // Дисковая информация провайдера
	Version() string // Версия провайдера
	Network() string // Сетевая информация провайдера
	Ram() string     // Информация о RAM провайдера
}

// NewController создает новый контроллер на основе указанной операционной системы.
// На данный момент поддерживается только Ubuntu.
func NewController(OS string) (API, error) {

	switch OS {
	case "Ubuntu":
		// Создаем и возвращаем контроллер для Ubuntu
		return &ubuntu.Ubuntu{}, nil
	default:
		// Если операционная система не известна, возвращаем ошибку
		return nil, errors.New("unknown controllers type")
	}
}
