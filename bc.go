package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main1() {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"}, // Укажите здесь адрес(а) вашего etcd-сервера
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Printf("Failed to create etcd client: %v\n", err)
		return
	}
	defer client.Close()

	// Создание пользователя "John" с паролем
	username := "John1"
	password := "password"

	_, err = client.UserAdd(context.Background(), username, password)
	if err != nil {
		fmt.Printf("Failed to add user: %v\n", err)
	}

	// Добавление роли с привилегиями
	_, err = client.RoleAdd(context.Background(), "readwriteRole")
	if err != nil {
		fmt.Printf("Failed to add role: %v\n", err)
	}

	_, err = client.RoleGrantPermission(context.Background(), "readwriteRole", "/home/",
		"", clientv3.PermissionType(clientv3.PermReadWrite))
	if err != nil {
		fmt.Printf("Failed to grant permission to role: %v\n", err)
		return
	}

	// Присвоение роли пользователю
	_, err = client.UserGrantRole(context.Background(), username, "readwriteRole")
	if err != nil {
		fmt.Printf("Failed to grant role to user: %v\n", err)
		return
	}

	// Используем присвоенные права для записи значения в ключ
	_, err = client.Put(context.Background(), "/home/", "Hello, RBAC!")
	if err != nil {
		fmt.Printf("Failed to put value: %v\n", err)
		return
	}

	// Используем присвоенные права для чтения значения из ключа
	resp, err := client.Get(context.Background(), "/home/")
	if err != nil {
		fmt.Printf("Failed to get value: %v\n", err)
		return
	}

	// Выводим значение на экран
	for _, ev := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", ev.Key, ev.Value)
	}
}
