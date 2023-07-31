package main

import (
	"BFG_auth/controllers"
	"BFG_auth/token_service"
	"github.com/pkg/errors"
	"time"
)

type Session struct {
	Username     string
	Token        string
	TokenManager token_service.TokenManager
	API          controllers.API
	Expiry       time.Time
}

func NewSession(username string, token string, tokenManager token_service.TokenManager, api controllers.API, expiry time.Time) *Session {
	return &Session{
		Username:     username,
		Token:        token,
		TokenManager: tokenManager,
		API:          api,
		Expiry:       expiry,
	}
}

func (s *Session) IsTokenValid() (bool, error) {
	_, err := s.TokenManager.ValidateToken(s.Token)
	if err != nil {
		return false, err
	}

	// Сверяем время истечения токена с текущим временем
	if time.Now().After(s.Expiry) {
		return false, nil
	}

	return true, nil
}

func (s *Session) ExecuteCommand(command string) (string, error) {
	// Проверяем валидность токена
	isValid, err := s.IsTokenValid()
	if err != nil || !isValid {
		return "", err
	}

	// Если токен валиден, выполняем команду через API
	controller, err := controllers.NewController(command)
	if err != nil {
		return "", err
	}

	switch command {
	case "Name":
		return controller.Name(), nil
	case "Time":
		return controller.Name(), nil
	case "Disk":
		return controller.Disk(), nil
	case "Version":
		return controller.Version(), nil
	case "Network":
		return controller.Network(), nil
	case "Ram":
		return controller.Ram(), nil
	default:
		return "", errors.New("no command in API")
	}
}
