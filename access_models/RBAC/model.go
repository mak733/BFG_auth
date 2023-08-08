// Package RBAC определяет основные операции над моделью CRUD, осуществляет хеширование и связь с БД.
package RBAC

// ServiceRBAC представляет сервис Role-Based Access Control для управления (CRUD) и хеширования иерархий пользователей.
// В качестве хранилища использует переданную БД.

import (
	"BFG_auth/access_models/accessTypes"
	"BFG_auth/repository"
	repoTypes "BFG_auth/repository/types"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

// Временные локальные хранилища для пользователей, ролей, групп и объектов
var (
	users   = make(map[string]*accessTypes.User)
	roles   = make(map[string]*accessTypes.Role)
	groups  = make(map[string]*accessTypes.Group)
	objects = make(map[string]*accessTypes.Object)
)

// Структура для хранения модели и связи с репозиторием
type ServiceRBAC struct {
	Repo repository.UserRepository
}

// readFromRepo запрашивает информацию о пользователе по его uid из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор пользователя
//
// Возвращает:
//   - error: ошибка, если заданный ключ не найден
//   - repoTypes.Key: ключ, соответствующий пользователю
//   - repoTypes.Value: значение, соответствующее пользователю
func (s ServiceRBAC) readFromRepo(uid accessTypes.Uid) (error, repoTypes.Key, repoTypes.Value) {
	kv, err := s.Repo.Read(repoTypes.Key(uid))
	if err != nil {
		return errors.New(fmt.Sprintf("No %s in repo", uid)), nil, nil
	}
	return nil, kv.Key, kv.Value
}

// writeToRepo сохраняет данные в репозиторий.
//
// Принимает:
//   - kv: пара ключ-значение для сохранения
//
// Возвращает:
//   - error: ошибка, если не удалось записать пару
func (s ServiceRBAC) writeToRepo(kv repoTypes.KV) error {
	err := s.Repo.Create(kv)
	if err != nil {
		return err
	}
	return nil
}

// CreateUser создает нового пользователя и сохраняет его в репозитории.
//
// Принимает:
//   - uid: уникальный идентификатор пользователя
//   - idp: идентификатор поставщика идентификации
//   - newRoles: список ролей, которые нужно присвоить пользователю
//   - newGroups: список групп, к которым нужно присоединить пользователя
//
// Возвращает:
//   - *accessTypes.User: созданный пользователь
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) CreateUser(uid, idp string, newRoles, newGroups []string) (*accessTypes.User, error) {
	// Инициализация нового пользователя
	user := accessTypes.User{
		Uid: accessTypes.Uid(uid),
		IdP: accessTypes.IdP(idp),
	}

	// Назначение ролей пользователю
	user.IdRoles = make([]accessTypes.IdRole, len(newRoles))
	user.Roles = make(map[accessTypes.IdRole]*accessTypes.Role, len(newRoles))
	for i, role := range newRoles {
		readRole, err := s.ReadRole(accessTypes.Uid(role))
		if err != nil {
			return nil, err
		}
		user.IdRoles[i] = readRole.Id
		user.Roles[readRole.Id] = readRole
	}

	// Назначение групп пользователю
	user.IdGroups = make([]accessTypes.IdGroup, len(newGroups))
	user.Groups = make(map[accessTypes.IdGroup]*accessTypes.Group, len(newGroups))
	for i, group := range newGroups {
		readGroup, err := s.ReadGroup(accessTypes.Uid(group))
		if err != nil {
			return nil, err
		}
		user.IdGroups[i] = readGroup.Id
		user.Groups[readGroup.Id] = readGroup
	}

	// Конвертация пользователя в JSON и запись в репозиторий
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonUser})
	if err != nil {
		return nil, err
	}

	// Добавление пользователя в локальное хранилище
	users[string(user.Uid)] = &user
	return &user, nil
}

// CreateRole создает новую роль и сохраняет ее в репозитории.
//
// Принимает:
//   - id: уникальный идентификатор роли
//   - newPermissions: разрешения, которые нужно присвоить роли
//
// Возвращает:
//   - *accessTypes.Role: созданная роль
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) CreateRole(id string,
	newPermissions map[accessTypes.IdObject]map[accessTypes.PermissionEnum]bool) (*accessTypes.Role, error) {
	// TODO: Проверка правильности newPermissions
	idObjects := make([]accessTypes.IdObject, len(newPermissions))
	i := 0
	for idObject := range newPermissions {
		idObjects[i] = idObject
		i++
	}

	// Инициализация новой роли
	role := accessTypes.Role{
		Id:        accessTypes.IdRole(id),
		IdObjects: idObjects,
	}
	role.Permissions = newPermissions

	// Конвертация роли в JSON и запись в репозиторий
	jsonRole, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}
	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(id), Value: jsonRole})
	if err != nil {
		return nil, err
	}
	// Добавление роли в локальное хранилище
	roles[string(role.Id)] = &role
	return &role, nil
}

// CreateGroup создает новую группу и сохраняет ее в репозитории.
//
// Принимает:
//   - id: уникальный идентификатор группы
//   - newRoles: список ролей, которые нужно присвоить группе
//
// Возвращает:
//   - *accessTypes.Group: созданная группа
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) CreateGroup(id string, newRoles []string) (*accessTypes.Group, error) {
	// Инициализация новой группы
	group := accessTypes.Group{
		Id: accessTypes.IdGroup(id),
	}
	group.IdRoles = make([]accessTypes.IdRole, len(newRoles))
	group.Roles = make(map[accessTypes.IdRole]*accessTypes.Role, len(newRoles))
	for i, role := range newRoles {
		readRole, err := s.ReadRole(accessTypes.Uid(role))
		if err != nil {
			return nil, err
		}
		group.IdRoles[i] = readRole.Id
		group.Roles[readRole.Id] = readRole
	}

	// Конвертация группы в JSON и запись в репозиторий
	jsonGroup, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}
	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(id), Value: jsonGroup})
	if err != nil {
		return nil, err
	}
	// Добавление группы в локальное хранилище
	groups[string(group.Id)] = &group
	return &group, nil
}

// CreateObject создает новый объект и сохраняет его в репозитории.
//
// Принимает:
//   - id: уникальный идентификатор объекта
//
// Возвращает:
//   - *accessTypes.Object: созданный объект
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) CreateObject(id accessTypes.IdObject) (*accessTypes.Object, error) {
	// Инициализация нового объекта
	var object accessTypes.Object
	object.Id = id

	// Конвертация объекта в JSON и запись в репозиторий
	jsonObject, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	err = s.writeToRepo(repoTypes.KV{Key: repoTypes.Key(id), Value: jsonObject})
	if err != nil {
		return nil, err
	}
	// Добавление группы в локальное хранилище
	objects[string(object.Id)] = &object
	return &object, nil
}

// ReadUser считывает информацию о пользователе из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор пользователя
//
// Возвращает:
//   - *accessTypes.User: считанный пользователь
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) ReadUser(uid accessTypes.Uid) (*accessTypes.User, error) {
	var user accessTypes.User

	// Считывание информации о пользователе из репозитория
	err, id, attributes := s.readFromRepo(uid)
	if err != nil {
		return nil, err
	}

	// Десериализация полученных данных в структуру User
	err = json.Unmarshal(attributes, &user)
	if err != nil {
		return nil, err
	}

	// Проверка всех ролей, связанных с пользователем
	for _, role := range user.IdRoles {
		_, err := s.ReadRole(accessTypes.Uid(role))
		if err != nil {
			return nil, err
		}
	}

	// Проверка всех групп, связанных с пользователем
	for _, group := range user.IdGroups {
		_, err := s.ReadGroup(accessTypes.Uid(group))
		if err != nil {
			return nil, err
		}
	}

	// Вывод информации о считанном пользователе
	fmt.Printf("RBAC: read user %s: %s\n", id, string(attributes))
	user.Uid = accessTypes.Uid(id)
	users[string(id)] = &user

	return &user, nil
}

// ReadRole считывает информацию о роли из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор роли
//
// Возвращает:
//   - *accessTypes.Role: считанная роль
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) ReadRole(uid accessTypes.Uid) (*accessTypes.Role, error) {
	// Считывание информации о роли из репозитория
	err, id, attributes := s.readFromRepo(uid)

	var role accessTypes.Role
	if err != nil {
		return nil, err
	}

	// Десериализация полученных данных в структуру Role
	err = json.Unmarshal(attributes, &role)
	if err != nil {
		return nil, err
	}

	// Вывод информации о считанной роли
	fmt.Printf("RBAC: read role %s: %s\n", id, string(attributes))
	roles[string(id)] = &role
	return &role, nil
}

// ReadGroup считывает информацию о группе из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор группы
//
// Возвращает:
//   - *accessTypes.Group: считанная группа
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) ReadGroup(uid accessTypes.Uid) (*accessTypes.Group, error) {
	var group accessTypes.Group

	// Считывание информации о группе из репозитория
	err, id, attributes := s.readFromRepo(uid)
	if err != nil {
		return nil, err
	}

	// Десериализация полученных данных в структуру Group
	err = json.Unmarshal(attributes, &group)
	if err != nil {
		return nil, err
	}

	// Вывод информации о считанной группе
	fmt.Printf("RBAC: read group %s: %s\n", id, string(attributes))
	groups[string(id)] = &group
	return &group, nil
}

// ReadObject считывает информацию об объекте из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор объекта
//
// Возвращает:
//   - *accessTypes.Object: считанный объект
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) ReadObject(uid accessTypes.Uid) (*accessTypes.Object, error) {
	// Считывание информации об объекте из репозитория
	err, id, attributes := s.readFromRepo(uid)

	var object accessTypes.Object
	if err != nil {
		return nil, err
	}

	// Десериализация полученных данных в структуру Object
	err = json.Unmarshal(attributes, &object)
	if err != nil {
		return nil, err
	}

	// Вывод информации о считанном объекте
	fmt.Printf("RBAC: read %s: %s\n", id, string(attributes))
	objects[string(id)] = &object
	return &object, nil
}

// UpdateUser обновляет информацию о пользователе в репозитории.
//
// Принимает:
//   - uid: уникальный идентификатор пользователя
//   - newUser: новые данные пользователя
//
// Возвращает:
//   - bool: true, если обновление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) UpdateUser(uid accessTypes.Uid, newUser *accessTypes.User) (bool, error) {
	// Обновляем пользователя в репозитории
	jsonUser, err := json.Marshal(newUser)
	if err != nil {
		return false, err
	}

	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(newUser.Uid), Value: jsonUser})
	if err != nil {
		return false, err
	}
	// Обновляем пользователя в локальной мапе
	users[string(uid)] = newUser
	return true, nil
}

// UpdateRole обновляет информацию о роли в репозитории.
//
// Принимает:
//   - uid: уникальный идентификатор роли
//   - newRole: новые данные роли
//
// Возвращает:
//   - bool: true, если обновление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) UpdateRole(uid accessTypes.Uid, newRole *accessTypes.Role) (bool, error) {
	// Конвертация новых данных роли в формат JSON
	jsonRole, err := json.Marshal(newRole)
	if err != nil {
		return false, err
	}

	// Обновление данных роли в репозитории
	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonRole})
	if err != nil {
		return false, err
	}

	// Обновление данных роли в локальной карте
	roles[string(uid)] = newRole
	return true, nil
}

// UpdateGroup обновляет информацию о группе в репозитории.
//
// Принимает:
//   - uid: уникальный идентификатор группы
//   - newGroup: новые данные группы
//
// Возвращает:
//   - bool: true, если обновление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) UpdateGroup(uid accessTypes.Uid, newGroup *accessTypes.Group) (bool, error) {
	// Конвертация новых данных группы в формат JSON
	jsonGroup, err := json.Marshal(newGroup)
	if err != nil {
		return false, err
	}

	// Обновление данных группы в репозитории
	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonGroup})
	if err != nil {
		return false, err
	}

	// Обновление данных группы в локальной карте
	groups[string(uid)] = newGroup
	return true, nil
}

// UpdateObject обновляет информацию об объекте в репозитории.
//
// Принимает:
//   - uid: уникальный идентификатор объекта
//   - newObject: новые данные объекта
//
// Возвращает:
//   - bool: true, если обновление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) UpdateObject(uid accessTypes.Uid, newObject *accessTypes.Object) (bool, error) {
	// Конвертация новых данных объекта в формат JSON
	jsonObject, err := json.Marshal(newObject)
	if err != nil {
		return false, err
	}

	// Обновление данных объекта в репозитории
	err = s.Repo.Update(repoTypes.KV{Key: repoTypes.Key(uid), Value: jsonObject})
	if err != nil {
		return false, err
	}

	// Обновление данных объекта в локальной карте
	objects[string(uid)] = newObject
	return true, nil
}

// DeleteUser удаляет пользователя из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор пользователя
//
// Возвращает:
//   - bool: true, если удаление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) DeleteUser(uid accessTypes.Uid) (bool, error) {
	// Удаление данных пользователя из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}

	// Удаление данных пользователя из локальной карте
	delete(users, string(uid))
	return true, nil
}

// DeleteRole удаляет роль из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор роли
//
// Возвращает:
//   - bool: true, если удаление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) DeleteRole(uid accessTypes.Uid) (bool, error) {
	// Удаление данных роли из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}

	// Удаление данных роли из локальной карте
	delete(roles, string(uid))
	return true, nil
}

// DeleteGroup удаляет группу из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор группы
//
// Возвращает:
//   - bool: true, если удаление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) DeleteGroup(uid accessTypes.Uid) (bool, error) {
	// Удаление данных группы из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}

	// Удаление данных группы из локальной карте
	delete(groups, string(uid))
	return true, nil
}

// DeleteObject удаляет объект из репозитория.
//
// Принимает:
//   - uid: уникальный идентификатор объекта
//
// Возвращает:
//   - bool: true, если удаление прошло успешно, иначе false
//   - error: ошибка, если таковая имеется
func (s ServiceRBAC) DeleteObject(uid accessTypes.Uid) (bool, error) {
	// Удаление данных объекта из репозитории
	err := s.Repo.Delete(repoTypes.Key(uid))
	if err != nil {
		return false, err
	}

	// Удаление данных объекта из локальной карте
	delete(objects, string(uid))
	return true, nil
}
