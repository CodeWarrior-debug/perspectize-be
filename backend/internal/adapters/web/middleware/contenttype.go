package middleware

import (
	"net/http"
	"strings"
)

// ContentTypeValidation rejects non-JSON POST requests to prevent CSRF.
// Browsers cannot send application/json via form submission, making CSRF
// impossible without CSRF tokens.
func ContentTypeValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ct := r.Header.Get("Content-Type")
			// Allow only application/json (with optional charset)
			if !strings.HasPrefix(ct, "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnsupportedMediaType)
				w.Write([]byte(`{"errors":[{"message":"Content-Type must be application/json"}]}`))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
