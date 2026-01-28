package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func RateLimit(requestsPerWindow int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.LimitByIP(requestsPerWindow, window)
}
