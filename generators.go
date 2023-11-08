package main

import (
	"log"
	"math/rand"
	"time"
)

func genereateUniqueAccount() int64 {
	timestamp := time.Now().UnixNano() / int64(time.Nanosecond) * 1000

	random := int64(rand.Intn(1000))

	uniqueID := timestamp + random

	return int64(uniqueID)
}

func generateSeeds(store Storage) {
	seedAccount(store, "admin@mail.com", "adminpassword", "admin", "admin", "1234567890")
}

func seedAccount(store Storage, email, pw, fname, lname, pNumber string) *User {
	user, err := NewAdminAccount(email, pw, fname, lname, pNumber)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateUser(user, &Account{}); err != nil {
		log.Fatal(err)
	}

	return user
}

func seedTransaction(store Storage, from, to int, description string, transactionType int) {
}
