package main

import (
	"BFG_auth/session"
	"BFG_auth/view"
	"fmt"
)

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

	abstractView, err := view.NewView("http")
	if err != nil {
		fmt.Printf("Error create abstractView %s: %v\n", "http", err)
		return
	}

	err = abstractView.StartServer(":8080", sessionManager)
	if err != nil {
		fmt.Printf("Error start abstractView %s: %v\n", "http", err)
		return
	}
}
