package main

import (
	"BFG_auth/session"

	"BFG_auth/view"
	view_types "BFG_auth/view/types"
	"fmt"
)

func handleUserSession(session session.UserSession, command, token string) {

	// здесь управление сессией пользователя, возможность вызывать API и т.д.
}

func main() {

	sessionManager, err := session.NewSessionManager(
		"Ubuntu",
		"RBAC",
		"etcd",
		"JWT",
	)
	if err != nil {
		fmt.Printf("Error to creaate session manager %s\n", err)
	}

	//  0. Создаём вьюху для обзения с юхверем
	view, err := view.NewView("http")
	if err != nil {
		fmt.Printf("Error create view %s: %v\n", "http", err)
		return
	}

	// Create a channel for communication between view and main
	authUserChannel := make(chan view_types.Credentials, 100)
	go view.StartServer(&authUserChannel, sessionManager)

	for {
		// Wait for credentials from the view
		credentials := <-authUserChannel
		username := credentials.Username
		password := credentials.Password
		fmt.Printf("New user %s %s\n", username, password)

		//чекнуть что сессия открыта для юзверя

		//session, err := authenticate(api, tokenManager, username, password)

		if err != nil {
			fmt.Printf("Error to authenticate user %s: %v\n", username, err)
			continue
		}

		//fmt.Printf("Token %s\n", session)

		//err := autorisation(session, repo, model)
		//credentials := <-authUserChannel
		//handleUserSession(session, apiChanel)
	}
}

/*
func autorisation(session UserSession, repo repository.UserRepository,
	model access_models.NewAccessControlModel) error {
	//идем в репо ищем юзера
	kv, err := repo.Read(types.Key(session.Username))

	if err != nil {
		return UserSession{IsAuthenticated: false}, err
	}

	//идем в модельку и ищем юзверя
	//model.CreateUser(access_types.Uid(kv.Key), kv.Value)
	user, err := model.ReadUser(access_types.Uid(kv.Key))
	if err != nil {
		return UserSession{IsAuthenticated: false}, err
	}
	fmt.Printf("%+v", user)

	return nil
}
*/
