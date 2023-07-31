package main

import (
	"BFG_auth/access_models"
	access_types "BFG_auth/access_models/types"
	"BFG_auth/controllers"
	"BFG_auth/identity_providers"
	"BFG_auth/repository"
	repo_types "BFG_auth/repository/types"
	"BFG_auth/token_service"
	"BFG_auth/view"
	view_types "BFG_auth/view/types"
	"fmt"
	"time"
)

func handleUserSession(api controllers.API, user, token string) {
	// здесь управление сессией пользователя, возможность вызывать API и т.д.
}

func main() {
	//  0. Создаём вьюху для обзения с юхверем
	viewName := "http"
	view, err := view.NewView(viewName)
	if err != nil {
		fmt.Printf("Error create view %s: %v\n", viewName, err)
		return
	}

	os := "Ubuntu"
	api, err := controllers.NewController(os)
	if err != nil {
		fmt.Printf("Error create API for os %s: %v\n", os, err)
		return
	}

	modelName := "RBAC"
	model, err := access_models.CreateAccessControlModel(modelName)
	if err != nil {
		fmt.Printf("Error make model %s\n", modelName)
		return
	}

	repoName := "etcd"
	repo, err := repository.NewRepository(repoName)
	if err != nil {
		fmt.Printf("Error make repository %s\n", repoName)
		return
	}

	// Create a channel for communication between view and main
	authUserChannel := make(chan view_types.Credentials, 100)
	go view.StartServer(&authUserChannel)

	for {
		// Wait for credentials from the view
		credentials := <-authUserChannel
		username := credentials.Username
		password := credentials.Password
		fmt.Printf("New user %s %s\n", username, password)
		token, err := authenticate(repo, model, api, username, password)
		if err != nil {
			fmt.Printf("Error to authenticate user %s: %v\n", username, err)
			return
		}
		fmt.Printf("Token %s\n", token)

		go handleUserSession(api, username, password)
	}
}

func authenticate(repo repository.UserRepository, model access_models.AccessControl,
	api controllers.API, username, password string) (string, error) {
	//идем в репо ищем юзера
	kv, err := repo.Read(repo_types.Key(username))

	if err != nil {
		return "", err
	}
	//идем в модельку и ищем юзверя

	//model.CreateUser(access_types.Uid(kv.Key), kv.Value)
	model.ReadUser(access_types.Uid(username))
	if err != nil {
		fmt.Println("Error create user:", err)
		return "", err
	}

	user, err := model.ReadUser(access_types.Uid(kv.Key))
	if err != nil {
		fmt.Println("Error read user:", err)
		return "", err
	}
	fmt.Printf("%+v", user)

	//	4. Проводим с помощью identity_providers аутентификацию
	IdP, err := identity_providers.NewIdp(string(user.IdP))
	isAuthenticate, err := IdP.Authenticate(string(user.Uid), password)

	if err != nil {
		fmt.Println("Error during authentication:", err)
		return "", err
	}
	fmt.Printf("User %s with password %s is %d\n", username, password, isAuthenticate)

	// 5. Если все хорошо то генерим токен

	if !isAuthenticate {
		fmt.Println("Not authenticated")
		return "", err
	}

	tokenManager, err := token_service.CreateTokenManager("JWT")
	token, err := tokenManager.GenerateToken(string(user.Uid))

	if err != nil {
		fmt.Println("Error during get token:", err)
		return "", err
	}
	fmt.Printf("User %s get new token for a %t\n", username, 24*time.Hour)

	return token, nil
}
