package main

import (
	"flag"
	"fmt"
	"log"
)

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
    // acc, err := 
}

func seedAccounts(s Storage) {
	seedAccount(s, "admin@mail.com", "adminpassword", "admin", "admin", "1234567890")
	seedAccount(s, "admin2@mail.com", "adminpassword", "admin2", "admin2", "1234567890")
}

func main() {
	seed := flag.Bool("seed", false, "seed the db with admin")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding the database")
		seedAccounts(store)
	}

	server := NewAPIServer(":3030", store)
	server.Run()
}
