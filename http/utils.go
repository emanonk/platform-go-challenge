package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		log.Printf("error writing json: %v", err)
	}
}

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m ...Middleware) http.Handler {
	// apply in reverse so Chain(h, a, b) = a(b(h))
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

// LoggingMiddleware logs request start/end to help visualize middleware order.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("http: start method=%s path=%s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("http: end method=%s path=%s dur=%s", r.Method, r.URL.Path, time.Since(start))
	})
}
