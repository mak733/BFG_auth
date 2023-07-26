package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type User struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func crateUser(ctx context.Context, client *clientv3.Client, username, password string) error {
	// Создание пользователя с паролем
	_, err := client.UserAdd(ctx, username, password)
	if err != nil {
		return fmt.Errorf("failed to add user: %v", err)
	}

	return nil
}

func readUser(ctx context.Context, client *clientv3.Client, username string) (*User, error) {
	resp, err := client.UserGet(ctx, username)
	if err != nil {
		return nil, err
	}

	if len(resp.Roles) == 0 {
		return nil, fmt.Errorf("user %s not found", username)
	}

	user := User{Username: username, Roles: resp.Roles}

	return &user, nil
}

func updateUser(ctx context.Context, client *clientv3.Client, username, newPassword string) error {
	_, err := client.UserChangePassword(ctx, username, newPassword)
	if err != nil {
		return fmt.Errorf("failed to update user password: %v", err)
	}

	return nil
}

func deleteUser(ctx context.Context, client *clientv3.Client, username string) error {
	_, err := client.UserDelete(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

func main() {
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

	// Входные данные для пользователя
	username := "John"
	//password := "mypassword"
	//newPassword := "mynewpassword"

	// Добавление пользователя
	//err = crateUser(context.Background(), client, username, password)
	if err != nil {
		fmt.Printf("Failed to add user: %s: %v\n", username, err)
		return
	}

	fmt.Println("User added successfully.")

	// Чтение пользователя
	user, err := readUser(context.Background(), client, username)
	if err != nil {
		fmt.Printf("Failed to read user %s: %v\n", username, err)
		return
	}

	fmt.Println("User read successfully: %s: %s", user.Username, user.Roles)

	// Обновление пароля пользователя
	//err = updateUser(context.Background(), client, username, newPassword)
	if err != nil {
		fmt.Printf("Failed to update user password: %s: %v\n", username, err)
		return
	}

	fmt.Println("User password updated successfully.")

	//err = deleteUser(context.Background(), client, username)
	if err != nil {
		fmt.Printf("Failed to delete user %s: %v\n", username, err)
		return
	}

	fmt.Println("User password updated successfully.")
}
