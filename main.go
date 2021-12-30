package main

import (
	"log"

	sw "github.com/pvr1/gigs/go"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(router.Run(":8080"))
}
