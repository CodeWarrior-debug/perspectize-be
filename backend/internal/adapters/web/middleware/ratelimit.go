package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// GlobalRateLimit returns middleware limiting requests per IP globally.
func GlobalRateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	return httprate.LimitByIP(requestsPerMinute, time.Minute)
}

// GraphQLRateLimit returns a stricter rate limit for the GraphQL endpoint.
// Returns a GraphQL-formatted JSON error on limit exceeded.
func GraphQLRateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	return httprate.Limit(
		requestsPerMinute, time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"errors":[{"message":"rate limit exceeded"}]}`))
		}),
	)
}
