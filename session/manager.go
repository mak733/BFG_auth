package session

import (
	"BFG_auth/access_models"
	"BFG_auth/access_models/accessTypes"
	"BFG_auth/controllers"
	"BFG_auth/identity_providers"
	"BFG_auth/token_service"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

func (sm *Manager) GetUser(token string) (*accessTypes.User, error) {
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()

	session, ok := sm.users[token]
	if !ok {
		return nil, errors.New("no session for token")
	}

	return session, nil
}

type Manager struct {
	mtx   sync.RWMutex
	users map[string]*accessTypes.User

	api          controllers.API
	accessModel  access_models.AccessControl
	tokenManager token_service.TokenManager
}

func NewSessionManager(apiName, modelName,
	repoName, tokenManagerName string) (*Manager, error) {

	api, err := controllers.NewController(apiName)
	if err != nil {
		fmt.Printf("Error create API for os %s: %v\n", apiName, err)
		return nil, err
	}

	model, err := access_models.NewAccessControlModel(modelName, repoName)
	if err != nil {
		fmt.Printf("Error make model %s\n", modelName)
		return nil, err
	}

	tokenManager, err := token_service.GetTokenManager(tokenManagerName)
	if err != nil {
		fmt.Printf("Error make token manager %s\n", tokenManagerName)
	}

	return &Manager{
		users:        make(map[string]*accessTypes.User),
		api:          api,
		accessModel:  model,
		tokenManager: tokenManager,
	}, nil
}

func (sm *Manager) CloseSession(token string) error {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	_, ok := sm.users[token]
	if !ok {
		return errors.New("No user in model")
	}

	delete(sm.users, token)
	return nil
}

func (sm *Manager) Authenticate(username, password string) (string, error) {

	idp := "ldap"
	//читаем пользака из etcd для определения idp
	/*user, err := sm.model.ReadUser(accessTypes.Uid(username))
	if user.IdP == "" {
		idp = "ldap"
	}*/
	///а что если айдипи из логина не равен айдипи из репо?
	if strings.Contains(username, "@") {
		parts := strings.Split(username, "@")
		username = parts[0]
		idp = parts[1]
	}

	//	4. Проводим с помощью identity_providers аутентификацию
	IdP, err := identity_providers.NewIdp(idp)
	if err != nil {
		return "", err
	}

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

	return token, nil
}

func (sm *Manager) Authorize(username, token string) error {
	//если уже авторизован то ретерн без ошибки
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()
	sessionUser, _ := sm.users[token]

	if (sessionUser != nil) && (string(sessionUser.Uid) != username) {
		return errors.New("Duplicate token! Session with given token already exists for another user.")
	}

	//добавим пользака в модель
	user, err := sm.accessModel.ReadUser(accessTypes.Uid(username))
	if err != nil {
		return err
	}
	//open sessionUser
	sm.users[token] = user

	return nil
}

func (sm *Manager) ValidateToken(token string) (bool, error) {
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()
	sessionUser, err := sm.GetUser(token)
	if err != nil {
		return false, errors.New("no sessionUser for token")
	}

	return sm.tokenManager.ValidateToken(string(sessionUser.Uid), token)
}

func (s *Manager) ExecuteCommand(user *accessTypes.User, command string) (string, error) {
	if user == nil {
		return "", errors.New("No user session")
	}

	if !stringInSlice(command, user.Objects) {
		return "", errors.New(fmt.Sprintf("No object %s for user %s", command, user.Uid))
	}

	switch command {
	case "Name":
		return s.api.Name(), nil
	case "Time":
		return s.api.Time(), nil
	case "Disk":
		return s.api.Disk(), nil
	case "Version":
		return s.api.Version(), nil
	case "Network":
		return s.api.Network(), nil
	case "Ram":
		return s.api.Ram(), nil
	default:
		return "", errors.New(fmt.Sprintf("no command in API \"%s\"", command))
	}
}
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
