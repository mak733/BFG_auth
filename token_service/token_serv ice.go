// Package token_service предоставляет функциональность для работы с токенами.
package token_service

import (
	"BFG_auth/token_service/JWT"
	"errors"
	"fmt"
	"sync"
)

// TokenManager определяет интерфейс для управления токенами.
type TokenManager interface {
	// GenerateToken генерирует новый токен для указанного имени пользователя.
	GenerateToken(username string) (string, error)
	// ValidateToken проверяет действительность токена для указанного имени пользователя.
	ValidateToken(username, token string) (bool, error)
}

var (
	mu            sync.Mutex                      // mu используется для обеспечения безопасности при использовании в многопоточной среде.
	tokenManagers = make(map[string]TokenManager) // tokenManagers содержит кеш для экземпляров TokenManager.
)

// GetTokenManager возвращает экземпляр менеджера токенов по его имени.
// Принимает:
// - managerName: имя менеджера токенов (например, "JWT")
// Возвращает:
// - TokenManager: интерфейс для работы с токенами
// - error: ошибка, если таковая имеется
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
