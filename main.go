package main

import (
	"log"
)

func main() {
    store, err := NewPostgresStore()
    if err != nil {
         log.Fatal(err)
    }

    if err := store.Init(); err != nil {
        log.Fatal(err)
    }
    server := NewAPIServer(":3030", store)
    server.Run()
}