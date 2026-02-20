package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func CORS(frontendURL string, demoFrontendURL string) func(http.Handler) http.Handler {
	allowedOrigins := []string{frontendURL}
	if frontendURL != "http://localhost:3000" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000")
	}
	if demoFrontendURL != "" {
		allowedOrigins = append(allowedOrigins, demoFrontendURL)
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
