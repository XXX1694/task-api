package middleware

import (
	"encoding/json"
	"net/http"
)

var validAPIKeys = map[string]bool{
	"secret12345": true,
	"admin-key":   true,
	"test-key":    true,
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")

		if !validAPIKeys[apiKey] {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
