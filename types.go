package main

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	Email       string `json:"email"`
	Password    string `json:"-"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Account struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	PhoneNumber   string    `json:"phoneNumber"`
	AccountNumber int64     `json:"accountNumber"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"createdAt"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

func NewAccount(email, password, firstName, lastName, phoneNumber string) (*Account, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println("Unable to hash password")
		return nil, err
	}
	return &Account{
		Email:         email,
		Password:      hashedPassword,
		FirstName:     firstName,
		LastName:      lastName,
		PhoneNumber:   phoneNumber,
		AccountNumber: int64(rand.Intn(1000000)),
		CreatedAt:     time.Now().UTC(),
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPW), nil
}
