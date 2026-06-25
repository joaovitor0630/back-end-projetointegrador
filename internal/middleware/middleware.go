package middleware

import (
	"log"
	"net/http"
	"time"
)

// allowedOrigins is the whitelist of origins permitted for CORS.
// Add production URLs here when deploying.
var allowedOrigins = map[string]bool{
	"http://localhost:5173":  true,
	"https://localhost:5173": true,
	"http://localhost:3000":  true,
	"http://localhost:8080":  true,
	// Produção no Render
	"https://front-end-projetointegrador-mjzn.onrender.com": true,
}

// CORS middleware adds cross-origin headers with origin whitelist
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Vary", "Origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// SecurityHeaders middleware adds standard security headers
func SecurityHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(self), microphone=(), geolocation=()")
		next(w, r)
	}
}

// MaxBodySize middleware limits request body size (25 MB)
func MaxBodySize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const maxSize = 25 << 20 // 25 MB
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		next(w, r)
	}
}

// Logger middleware logs the HTTP method, URI, and time taken
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %v", r.Method, r.RequestURI, time.Since(start))
	}
}

// Wrap is a utility to apply multiple middlewares to a HandlerFunc
func Wrap(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
