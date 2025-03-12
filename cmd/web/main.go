package main

import (
	"log"
	"net/http"
)

func main() {
	mux := routes()

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
