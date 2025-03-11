package main

import (
	"log"
	"net/http"

	"github.com/linyows/ogenexample/api"
)

func main() {
	srv, err := api.Server()
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}
