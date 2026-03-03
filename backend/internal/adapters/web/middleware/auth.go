package middleware

import (
	"context"
	"net/http"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
)

// contextKey is an unexported type for context keys to prevent collisions.
type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware extracts a JWT from the "auth_token" httpOnly cookie,
// validates it, and stores the user claims in the request context.
// Invalid or missing tokens are treated as unauthenticated (request continues).
func AuthMiddleware(authService portservices.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				// No cookie — unauthenticated request, allowed to pass through
				next.ServeHTTP(w, r)
				return
			}

			claims, err := authService.ValidateToken(cookie.Value)
			if err != nil {
				// Invalid token — treat as unauthenticated, no error response
				next.ServeHTTP(w, r)
				return
			}

			// Store claims in context for downstream handlers
			ctx := context.WithValue(r.Context(), userContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ForContext extracts the authenticated user ID from the request context.
// Returns (userID, true) if authenticated, or (0, false) if not.
func ForContext(ctx context.Context) (int, bool) {
	claims, ok := ctx.Value(userContextKey).(*domain.Claims)
	if !ok || claims == nil {
		return 0, false
	}
	return claims.UserID, true
}
