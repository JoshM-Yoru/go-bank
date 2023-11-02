package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, email, pw, fname, lname, pNumber string) *Account {
	acc, err := NewAccount(email, pw, fname, lname, pNumber)
	if err != nil {
		log.Fatal(err)
	}

    if err := store.CreateAccount(acc); err != nil {
        log.Fatal(err)
    }

    fmt.Println("new account: ", acc.AccountNumber)

    return acc
}

func seedAccounts(s Storage) {
    seedAccount(s, "dmitry@mail", "passwordinrussian", "Dimitry", "Bivol", "1234567890")
}

func main() {
    seed := flag.Bool("seed", false, "seed the db")
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
