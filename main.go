package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
	log.Printf("Server started at %s", "8080")
}