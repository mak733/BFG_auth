package main

import (
	"BFG_auth/IdP"
	"BFG_auth/access_models"
	"BFG_auth/access_models/types"
	"context"
	"encoding/binary"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {
	//	1. Прилетает запрос юзер/пароль
	//  ...
	//	2. Ищем юезера в БД, цепляем какой IdP для него установлен
	username := "username"
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

	resp, err := client.Get(context.Background(), username)
	if err != nil {
		log.Fatalf("Failed to get the response: %v", err)
	}

	if resp.Count == 0 {
		log.Fatalf("Failed to get user: %s", username)
		return
	}

	body := resp.Kvs[0]
	fmt.Printf("Get value by key %s: %s\n", body.Key, body.Value)

	model, err := access_models.CreateAccessControlModel("RBAC")

	model.CreateUser(types.IdUser(binary.BigEndian.Uint64(body.Key)), string(body.Key), body.Value)

	//	3. Проводим с помощью IdP аутентификацию
	password := "password"
	IdP, err := IdP.NewIdp("ldap")
	isAuthenticate, err := IdP.Authenticate("testuser", "newpassword")

	if err != nil {
		fmt.Println("Error during authentication:", err)
		return
	}
	fmt.Printf("User %s with password %s is %d\n", username, password, isAuthenticate)
	/*	4. Дальше плешем по ТЗ*/
}
