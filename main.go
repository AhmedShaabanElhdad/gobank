package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {

	seed := flag.Bool("seed", false, "this to provide seed account to db for testing ")
	flag.Parse()

	store, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", store)
	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding account to DB")
		seedAccounts(store)
	}

	server := NewApiServer(":3001", store)
	fmt.Printf("connect to the server : localhost:30001/")
	server.Run()
}

func seedAccount(store Storage, fname string, lname string, password string) *Account {
	account, err := NewAccount(fname, lname, password)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateAccount(account); err != nil {
		log.Fatal(err)
	}

	return account
}

func seedAccounts(store Storage) {
	seedAccount(store, "Ahmed", "Shaban", "12345678")
}
