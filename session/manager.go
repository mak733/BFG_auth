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

// Manager представляет собой менеджер сессий для управления пользователями.
type Manager struct {
	mtx   sync.RWMutex
	users map[string]*accessTypes.User

	api          controllers.API
	accessModel  access_models.AccessControl
	tokenManager token_service.TokenManager
}

// GetUser возвращает пользователя по его токену.
// Принимает:
// - token: токен пользователя
// Возвращает:
// - *accessTypes.User: информация о пользователе
// - error: ошибка, если таковая имеется
func (m *Manager) GetUser(token string) (*accessTypes.User, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	session, ok := m.users[token]
	if !ok {
		return nil, errors.New("no session for token")
	}

	return session, nil
}

// NewSessionManager создает и возвращает новый экземпляр менеджера сессий.
// Принимает:
// - apiName: имя API
// - modelName: имя модели доступа
// - repoName: имя репозитория
// - tokenManagerName: имя менеджера токенов
// Возвращает:
// - *Manager: экземпляр менеджера сессий
// - error: ошибка, если таковая имеется
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

	///TEST

	///create role for user
	var readPermission accessTypes.PermissionEnum
	readPermission = "Read"
	t := make(map[accessTypes.IdObject]map[accessTypes.PermissionEnum]bool)
	// Check if the map for "Time" exists; if not, create it
	if t["Time"] == nil {
		t["Time"] = make(map[accessTypes.PermissionEnum]bool)
	}
	// Now you can assign a value for the permission
	t["Time"][readPermission] = true
	_, err = model.CreateRole("R_Admin", t)

	///create role for group
	m := make(map[accessTypes.IdObject]map[accessTypes.PermissionEnum]bool)
	// Check if the map for "Time" exists; if not, create it
	if m["Name"] == nil {
		m["Name"] = make(map[accessTypes.PermissionEnum]bool)
	}
	// Now you can assign a value for the permission
	m["Name"][readPermission] = true
	_, err = model.CreateRole("R_AdminG", m)

	//create group
	_, err = model.CreateGroup("G_Admin", []string{"R_AdminG"})

	//create user
	user, err := model.CreateUser("Admin", "ldap", []string{"R_Admin"}, []string{"G_Admin"})
	if err != nil {
		return nil, err
	}
	//update user
	user.Uid = "admin"
	user.IdP = "ldap"
	_, err = model.UpdateUser(user.Uid, user)
	if err != nil {
		return nil, err
	}

	//end TEST

	return &Manager{
		users:        make(map[string]*accessTypes.User),
		api:          api,
		accessModel:  model,
		tokenManager: tokenManager,
	}, nil
}

// CloseSession закрывает сессию пользователя по его токену.
// Принимает:
// - token: токен пользователя
// Возвращает:
// - error: ошибка, если таковая имеется
func (m *Manager) CloseSession(token string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	_, ok := m.users[token]
	if !ok {
		return errors.New("No user in model")
	}

	delete(m.users, token)
	return nil
}

// Authenticate аутентифицирует пользователя по его имени и паролю.
// Принимает:
// - username: имя пользователя
// - password: пароль пользователя
// Возвращает:
// - string: токен аутентификации
// - error: ошибка, если таковая имеется
func (m *Manager) Authenticate(username, password string) (string, error) {

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
	fmt.Printf("User %s with password %s is %b\n",
		username, password, isAuthenticate)

	if !isAuthenticate {
		return "", err
	}

	token, err := m.tokenManager.GenerateToken(username)

	if err != nil {
		return "", err
	}
	fmt.Printf("User %s get new token for a %d\n", username, 24*time.Hour)

	return token, nil
}

// Authorize авторизует пользователя по его имени и токену.
// Принимает:
// - username: имя пользователя
// - token: токен аутентификации
// Возвращает:
// - error: ошибка, если таковая имеется
func (m *Manager) Authorize(username, token string) error {
	//если уже авторизован то ретерн без ошибки
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	sessionUser, _ := m.users[token]

	if (sessionUser != nil) && (string(sessionUser.Uid) != username) {
		return errors.New("Duplicate token! Session with given token already exists for another user.")
	}

	//добавим пользака в модель
	user, err := m.accessModel.ReadUser(accessTypes.Uid(username))
	if err != nil {
		return err
	}
	//open sessionUser
	m.users[token] = user

	return nil
}

// ValidateToken проверяет действительность токена для пользователя.
// Принимает:
// - token: токен аутентификации
// Возвращает:
// - bool: true, если токен действителен, иначе false
// - error: ошибка, если таковая имеется
func (m *Manager) ValidateToken(token string) (bool, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	sessionUser, err := m.GetUser(token)
	if err != nil {
		return false, errors.New("no sessionUser for token")
	}

	return m.tokenManager.ValidateToken(string(sessionUser.Uid), token)
}

// ExecuteCommand выполняет команду для пользователя.
// Принимает:
// - user: информация о пользователе
// - command: команда для выполнения
// Возвращает:
// - string: результат выполнения команды
// - error: ошибка, если таковая имеется
func (m *Manager) ExecuteCommand(user *accessTypes.User, command string) (string, error) {
	if user == nil {
		return "", errors.New("No user session")
	}

	if ok, err := user.CheckPermission(accessTypes.IdObject(command), "Read"); !ok {
		return "", err
	}

	switch command {
	case "Name":
		return m.api.Name(), nil
	case "Time":
		return m.api.Time(), nil
	case "Disk":
		return m.api.Disk(), nil
	case "Version":
		return m.api.Version(), nil
	case "Network":
		return m.api.Network(), nil
	case "Ram":
		return m.api.Ram(), nil
	default:
		return "", errors.New(fmt.Sprintf("no command in API \"%s\"", command))
	}
}
