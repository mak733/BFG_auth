package main

import (
	"BFG_auth/IdP"
	"BFG_auth/access_models"
	"BFG_auth/access_models/types"
	"BFG_auth/repository"
	types_repo "BFG_auth/repository/types"
	"fmt"
)

func main() {
	//	1. Прилетает запрос юзер/пароль
	//  ...
	//	2. Ищем юезера в БД, цепляем какой IdP для него установлен
	username := "username"

	//идем в репо за юзверем

	repo, err := repository.NewRepository("etcd")
	kv, err := repo.Read(types_repo.Key(username))

	//идем в модельку и добавляем юзверя

	model, err := access_models.CreateAccessControlModel("RBAC")

	model.CreateUser(types.Uid(kv.Key), kv.Value)
	if err != nil {
		fmt.Println("Error create user:", err)
		return
	}

	user, err := model.ReadUser(types.Uid(kv.Key))
	if err != nil {
		fmt.Println("Error read user:", err)
		return
	}
	fmt.Printf("%+v", user)

	//	3. Проводим с помощью IdP аутентификацию
	password := "password"
	IdP, err := IdP.NewIdp(string(user.IdP))
	isAuthenticate, err := IdP.Authenticate("testuser", "newpassword")

	if err != nil {
		fmt.Println("Error during authentication:", err)
		return
	}
	fmt.Printf("User %s with password %s is %d\n", username, password, isAuthenticate)
	/*	4. Дальше плешем по ТЗ*/
}
