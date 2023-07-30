package main

import (
	"BFG_auth/controllers"
	"BFG_auth/token_service"
	"time"
)

type Session struct {
	Username     string
	Token        string
	TokenManager token_service.TokenManager
	API          controllers.Controller
	Expiry       time.Time
}

func NewSession(username string, token string, tokenManager token_service.TokenManager, api controllers.Controller, expiry time.Time) *Session {
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
	/*response, err := s.API.CallAPI(command, s.Token)
	if err != nil {
		return "", err
	}*/

	return response, nil
}
