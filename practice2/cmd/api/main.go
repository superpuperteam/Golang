package main

import (
	"assignment2/internal/handlers"
	"assignment2/internal/middleware"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/user", middleware.AuthAndLog(http.HandlerFunc(handlers.UserHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("listening on :8080")
	log.Fatal(server.ListenAndServe())
}
