package http

import (
	"log"
	"net/http"
	"time"
)

// LogMiddleware log the incoming request
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)

		log.Printf("%s %s elapsed: %s", r.Method, r.URL, elapsed)
	})
}
