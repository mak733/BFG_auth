// Package main
// В результате выполнения проекта должно быть разработано приложение на Golang, реализующее методы API с применением
// RBAC и хранением данных в ETCD. Приложение должно успешно работать на хост системе Linux и удовлетворять всем
// поставленным требованиям.
package main

import (
	"BFG_auth/session"
	"BFG_auth/view"
	"fmt"
)

func main() {
	// Инициализация нового менеджера сессий с заданными конфигурациями.
	// В данном примере используется "Ubuntu" как ОС, "RBAC" для контроля доступа, "etcd" как хранилище и "JWT" для токенов.
	sessionManager, err := session.NewSessionManager(
		"Ubuntu",
		"RBAC",
		"etcd",
		"JWT",
	)
	if err != nil {
		fmt.Printf("Ошибка при создании менеджера сессий: %s\n", err)
	}

	// Инициализация нового абстрактного представления.
	// В этом примере настраивается представление на основе HTTP.
	abstractView, err := view.NewView("http")
	if err != nil {
		fmt.Printf("Ошибка при создании представления %s: %v\n", "http", err)
		return
	}

	// Запуск HTTP-сервера на порту 8080 с использованием ранее инициализированного менеджера сессий.
	err = abstractView.StartServer(":8080", sessionManager)
	if err != nil {
		fmt.Printf("Ошибка при запуске представления %s: %v\n", "http", err)
		return
	}
}
