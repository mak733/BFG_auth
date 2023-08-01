package session

import (
	"BFG_auth/access_models"
	"BFG_auth/controllers"
	"BFG_auth/identity_providers"
	"BFG_auth/repository"
	"BFG_auth/token_service"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	mtx      sync.RWMutex
	sessions map[string]*UserSession

	api          controllers.API
	model        access_models.AccessControl
	repo         repository.UserRepository
	tokenManager token_service.TokenManager
}

func NewSessionManager(apiName, modelName,
	repoName, tokenManagerName string) (*Manager, error) {

	api, err := controllers.NewController(apiName)
	if err != nil {
		fmt.Printf("Error create API for os %s: %v\n", apiName, err)
		return nil, err
	}

	model, err := access_models.NewAccessControlModel(modelName)
	if err != nil {
		fmt.Printf("Error make model %s\n", modelName)
		return nil, err
	}

	repo, err := repository.NewRepository(repoName)
	if err != nil {
		fmt.Printf("Error make repository %s\n", repoName)
		return nil, err
	}

	tokenManager, err := token_service.GetTokenManager(tokenManagerName)
	if err != nil {
		fmt.Printf("Error make token manager %s\n", tokenManagerName)
	}

	return &Manager{
		sessions:     make(map[string]*UserSession),
		api:          api,
		model:        model,
		repo:         repo,
		tokenManager: tokenManager,
	}, nil
}

func (sm *Manager) CreateSession(username, token string) error {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	//вся логика из мейна

	//идем в репо ищем юзера
	/*kv, err := sm.repo.Read(types.Key(username))

	if err != nil {
		return err
	}

	//идем в модельку и ищем юзверя
	//model.CreateUser(access_types.Uid(kv.Key), kv.Value)
	user, err := sm.model.ReadUser(types2.Uid(kv.Key))
	if err != nil {
		return err
	}
	fmt.Printf("%+v", user)
	*/
	fmt.Printf("create %s\n", token)
	sm.sessions[token] = &UserSession{
		Username:        username,
		Token:           token,
		IsAuthenticated: true,
		TokenManager:    sm.tokenManager,
		API:             sm.api,
		repo:            sm.repo,
		accessModel:     sm.model,
	}
	return nil
}

func (sm *Manager) Authenticate(username, password string) (string, error) {
	idp := "ldap"

	if strings.Contains(username, "@") {
		parts := strings.Split(username, "@")
		username = parts[0]
		idp = parts[1]
	}

	//	4. Проводим с помощью identity_providers аутентификацию
	IdP, err := identity_providers.NewIdp(idp)
	isAuthenticate, err := IdP.Authenticate(username, password)

	if err != nil {
		return "", err
	}
	fmt.Printf("User %s with password %s is %d\n",
		username, password, isAuthenticate)

	if !isAuthenticate {
		return "", err
	}

	token, err := sm.tokenManager.GenerateToken(username)

	if err != nil {
		return "", err
	}
	fmt.Printf("User %s get new token for a %t\n", username, 24*time.Hour)

	err = sm.CreateSession(username, token)
	if err != nil {
		return "", err
	}
	//добавить экспайред тайм
	return token, nil
}

func (sm *Manager) GetSession(token string) (*UserSession, error) {
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()

	session, ok := sm.sessions[token]
	if !ok {
		return nil, errors.New("no session for token")
	}

	return session, nil
}

func (sm *Manager) ValidateToken(token string) (bool, error) {
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()
	session, err := sm.GetSession(token)
	if err != nil {
		return false, errors.New("no session for token")
	}

	return session.TokenManager.ValidateToken(session.Username, token)
}

type UserSession struct {
	Username        string
	Token           string
	IsAuthenticated bool
	TokenManager    token_service.TokenManager
	API             controllers.API
	repo            repository.UserRepository
	accessModel     access_models.AccessControl
}

func NewSession(username string, token string,
	tokenManager *token_service.TokenManager,
	api *controllers.API) UserSession {
	return UserSession{
		Username:     username,
		Token:        token,
		TokenManager: *tokenManager,
		API:          *api,
	}

}

func (s *UserSession) IsTokenValid() (bool, error) {

	//	manager := s.TokenManager

	return s.IsAuthenticated, nil
}

func (s *UserSession) ExecuteCommand(command string) (string, error) {
	// Проверяем валидность токена
	//isValid, err := s.TokenManager
	//if err != nil || !isValid {
	//	return "", err
	//}
	//чекнуть рбак права для сессии
	controller := s.API
	// Если токен валиден, выполняем команду через API
	switch command {
	case "Name":
		return controller.Name(), nil
	case "Time":
		return controller.Time(), nil
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
