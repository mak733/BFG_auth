package main

import "BFG_auth/repository"

func main() {
	r, _ := repository.NewRepository("etcd")
	r.Delete(1)
}
