package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	Email       string `json:"email"`
	Password    string `json:"-"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	UserName    string `json:"userName"`
	PhoneNumber string `json:"phoneNumber"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccountNumber int64  `json:"accountNumber"`
	Token         string `json:"token"`
}

type Role int

const (
	Admin Role = iota
	Employee
	Customer
	Guest
)

type TransactionType int

const (
	Debit TransactionType = iota
	Credit
)

type Account struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	UserName      string    `json:"userName"`
	PhoneNumber   string    `json:"phoneNumber"`
	AccountNumber int64     `json:"accountNumber"`
	Balance       int64     `json:"balance"`
	CreatedAt     time.Time `json:"createdAt"`
	Role          Role      `json:"role"`
	IsActive      bool      `json:"isActive"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

func NewAdminAccount(email, password, firstName, lastName, phoneNumber string) (*Account, error) {
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
		UserName:      "$" + firstName + "." + lastName + "#" + strconv.Itoa(int(rand.Intn(9000)+1000)),
		PhoneNumber:   phoneNumber,
		AccountNumber: int64(rand.Intn(1000000)),
		CreatedAt:     time.Now().UTC(),
		Role:          Admin,
	}, nil
}

func NewEmployeeAccount(email, password, firstName, lastName, phoneNumber string) (*Account, error) {
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
		UserName:      "$" + firstName + "." + lastName + "#" + strconv.Itoa(int(rand.Intn(9000)+1000)),
		PhoneNumber:   phoneNumber,
		AccountNumber: int64(rand.Intn(1000000)),
		CreatedAt:     time.Now().UTC(),
		Role:          Employee,
		IsActive:      true,
	}, nil
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
		UserName:      "$" + firstName + "." + lastName + "#" + strconv.Itoa(int(rand.Intn(9000)+1000)),
		PhoneNumber:   phoneNumber,
		AccountNumber: int64(rand.Intn(1000000)),
		CreatedAt:     time.Now().UTC(),
		Role:          Customer,
		IsActive:      true,
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPW), nil
}

func (account *Account) validatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(pw)) == nil
}
