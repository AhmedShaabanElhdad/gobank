package main

import (
	"math/rand"
	"time"
)

type AccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	ID        int       `json:"id"`
	Number    int64     `json:"number"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Balance   int64     `json:"balance"`
	createdAt time.Time `json:"createdAt"`
}

func NewAccount(firstName string, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000000)),
		createdAt: time.Now().UTC(),
	}
}
