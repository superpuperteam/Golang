package middleware

import (
	"assignment2/internal/handlers"
	"log"
	"net/http"
)

const apiKey = "secret123"

func AuthAndLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		if r.Header.Get("X-API-Key") != apiKey {
			handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ErrorResponse{Error: "unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
