package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role int

const (
	Admin Role = iota + 1
	Employee
	Customer
)

type AccountType int

const (
	Checking AccountType = iota + 1
	Savings
)

type TransactionType int

const (
	Debit TransactionType = iota + 1
	Credit
	Transfer
)

type CreateUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"-"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	UserName    string `json:"userName"`
	PhoneNumber string `json:"phoneNumber"`
	ReferrerID  string `json:"referrerID"`
	Balance     int    `json:"balance"`
	Role        string `json:"role"`
	AccountType string `json:"accountType"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TransactionRequest struct {
	FromAccount     int `json:"fromAccount"`
	ToAccount       int `json:"toAccount"`
	Amount          int `json:"amount"`
	TransactionType int `json:"transactionType"`
}

type LoginResponse struct {
	UserName string `json:"userName"`
	Token    string `json:"token"`
}

type User struct {
	ID          int       `json:"user_id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	UserName    string    `json:"userName"`
	PhoneNumber string    `json:"phoneNumber"`
	ReferrerID  int       `json:"referrerID"`
	CreatedAt   time.Time `json:"createdAt"`
	LastLogin   time.Time `json:"lastLogin"`
	Role        Role      `json:"role"`
	IsActive    bool      `json:"isActive"`
}

type Account struct {
	ID              int         `json:"account_id"`
	UserID          int         `json:"user_id"`
	AccountNumber   int64       `json:"accountNumber"`
	Balance         int64       `json:"balance"`
	CreatedAt       time.Time   `json:"createdAt"`
	AccountType     AccountType `json:"accountType"`
	IsActiveAccount bool        `json:"isActiveAccount"`
}

type FullAccount struct {
	User     User
	Accounts []Account
}

func NewAdminAccount(email, password, firstName, lastName, phoneNumber string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println("Unable to hash password")
		return nil, err
	}
	return &User{
		Email:       email,
		Password:    hashedPassword,
		FirstName:   firstName,
		LastName:    lastName,
		UserName:    "$" + firstName + "." + lastName + "#" + strconv.Itoa(int(rand.Intn(9000)+1000)),
		PhoneNumber: "",
		CreatedAt:   time.Now().UTC(),
		LastLogin:   time.Now().UTC(),
		Role:        Admin,
		IsActive:    true,
	}, nil
}

func NewUserAccount(email, password, firstName, lastName, phoneNumber string, referrerID int, balance int64, role Role, accType AccountType) (*User, *Account, error) {

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println("Unable to hash password")
		return nil, nil, err
	}
	return &User{
			Email:       email,
			Password:    hashedPassword,
			FirstName:   firstName,
			LastName:    lastName,
			UserName:    "$" + firstName + "." + lastName + "#" + strconv.Itoa(int(rand.Intn(9000)+1000)),
			PhoneNumber: phoneNumber,
			ReferrerID:  referrerID,
			CreatedAt:   time.Now().UTC(),
			Role:        role,
			IsActive:    true,
		}, &Account{
			AccountNumber:   int64(rand.Intn(990000) + 100000),
			Balance:         balance,
			CreatedAt:       time.Now().UTC(),
			AccountType:     accType,
			IsActiveAccount: true,
		}, nil
}

func hashPassword(password string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPW), nil
}

func (user *User) validatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw)) == nil
}
