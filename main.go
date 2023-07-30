package main

import (
	"BFG_auth/access_models"
	"BFG_auth/access_models/types"
	"BFG_auth/controllers"
	"BFG_auth/identity_providers"
	"BFG_auth/repository"
	types2 "BFG_auth/repository/types"
	"BFG_auth/token_service"
	"BFG_auth/view"
	"fmt"
	"time"
)

func main() {
	//  0. Создаём вьюху для обзения с юхверем

	viewName := "http"
	view, err := view.NewView(viewName)
	if err != nil {
		fmt.Print("Error create view %s: %s\n", viewName, err)
		return
	}
	// Тут мы ждём прихода нового юзера, когда пришел лвим юхернейм и пароль от него
	//  ...
	//	2. Ищем юезера в БД, цепляем какой identity_providers для него установлен

	repo, err := repository.NewRepository("etcd")
	kv, err := repo.Read(types2.Key(username))

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

	//	4. Проводим с помощью identity_providers аутентификацию
	password := "password"
	IdP, err := identity_providers.NewIdp(string(user.IdP))
	isAuthenticate, err := IdP.Authenticate(string(user.Uid), password)

	if err != nil {
		fmt.Println("Error during authentication:", err)
		return
	}
	fmt.Printf("User %s with password %s is %d\n", username, password, isAuthenticate)

	// 5. Если все хорошо то генерим токен

	if !isAuthenticate {
		fmt.Println("Not authenticated")
		return
	}

	tomenManager, err := token_service.CreateTokenManager("JWT")
	token, err := tomenManager.GenerateToken(string(user.Uid))

	if err != nil {
		fmt.Println("Error during get token:", err)
		return
	}
	fmt.Printf("User %s get new token for a %t\n", username, 24*time.Hour)

	// 6. Создаём контроллер(АПИ).
	os := "Ubuntu"
	api, err := controllers.NewController(os)

	if err != nil {
		fmt.Println("Error during get API for %s:", os, err)
		return
	}
	fmt.Printf("Use API for %s\n", os)

	// 7. Тут мы открываем сесси в которой и будет происходить весь бизнес,
	//	живет она пока токен жив, работает в потоке

}
