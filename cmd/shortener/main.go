package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /", shortHandler)
	http.HandleFunc("GET /{id}", redirectHandler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
