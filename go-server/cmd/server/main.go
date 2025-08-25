package main

import (
	"log"
	"net/http"
	httpserver "wallet-backend/internal/http"
)

func main() {
	h := httpserver.New()
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Router()))
}
