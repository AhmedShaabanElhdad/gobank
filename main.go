package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", store)
	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":3000", store)
	server.Run()
}
