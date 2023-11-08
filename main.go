package main

import (
	"flag"
	"fmt"
	"log"
)


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
		generateSeeds(store)
	}

	server := NewAPIServer(":3030", store)
	server.Run()
}
