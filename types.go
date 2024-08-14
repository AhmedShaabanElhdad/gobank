package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type TransferRequest struct {
	ToAccount int     `json:"toAccount"`
	Amount    float64 `json:"amout"`
	Details   string  `json:"details"`
}

type Account struct {
	ID                int       `json:"id"`
	Number            int64     `json:"number"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	EncryptedPassword string    `json:"-"` // this mean ignore when return
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (a *Account) validatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password)) == nil
}

func NewAccount(firstName string, lastName string, password string) (*Account, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:                rand.Intn(10000),
		FirstName:         firstName,
		LastName:          lastName,
		EncryptedPassword: string(encryptedPass),
		Number:            int64(rand.Intn(10000000)),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
