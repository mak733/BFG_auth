package token_service

import (
	"BFG_auth/token_service/JWT"
	"errors"
	"fmt"
	"sync"
)

type TokenManager interface {
	GenerateToken(username string) (string, error)
	ValidateToken(username, token string) (bool, error)
}

var (
	mu            sync.Mutex                      // Для обеспечения безопасности при использовании в многопоточной среде
	tokenManagers = make(map[string]TokenManager) // Кеш для экземпляров TokenManager
)

func GetTokenManager(managerName string) (TokenManager, error) {
	mu.Lock()
	defer mu.Unlock()

	// Если менеджер уже существует, вернуть его
	if manager, exists := tokenManagers[managerName]; exists {
		return manager, nil
	}

	// В противном случае создать новый
	switch managerName {
	case "JWT":
		fmt.Printf("Get %s\n", managerName)
		manager := JWT.NewTokenService() // Предполагается, что у вас есть соответствующий код
		tokenManagers[managerName] = manager
		return manager, nil
	default:
		return nil, errors.New("unknown token manager type")
	}
}
