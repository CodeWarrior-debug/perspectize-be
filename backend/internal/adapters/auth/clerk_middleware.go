package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	repositories "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/repositories"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
)

// Middleware verifies Clerk Bearer tokens and resolves local users.
// Permissive: unauthenticated requests pass through for public queries.
func Middleware(userRepo repositories.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Wrap with Clerk's JWT verification middleware
		return clerkhttp.WithHeaderAuthorization()(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims, ok := clerk.SessionClaimsFromContext(r.Context())
				if !ok || claims == nil {
					// No valid session — pass through as unauthenticated
					next.ServeHTTP(w, r)
					return
				}

				// Resolve Clerk user ID to local user
				clerkUserID := claims.Subject
				user, err := userRepo.GetByClerkID(r.Context(), clerkUserID)
				if err != nil {
					if errors.Is(err, domain.ErrNotFound) {
						// On-demand creation: webhook may not have fired yet
						clerkUsr, fetchErr := clerkuser.Get(r.Context(), clerkUserID)
						if fetchErr != nil {
							slog.Warn("clerk user not found via API",
								"clerk_user_id", clerkUserID,
								"error", fetchErr,
							)
							next.ServeHTTP(w, r)
							return
						}

						// Extract username and email from Clerk profile
						username := ""
						if clerkUsr.Username != nil {
							username = *clerkUsr.Username
						}
						email := ""
						if len(clerkUsr.EmailAddresses) > 0 {
							for _, ea := range clerkUsr.EmailAddresses {
								if clerkUsr.PrimaryEmailAddressID != nil && ea.ID == *clerkUsr.PrimaryEmailAddressID {
									email = ea.EmailAddress
									break
								}
							}
							if email == "" {
								email = clerkUsr.EmailAddresses[0].EmailAddress
							}
						}
						if username == "" {
							// Fall back to email prefix
							for i, c := range email {
								if c == '@' {
									username = email[:i]
									break
								}
							}
							if username == "" {
								username = clerkUserID
							}
						}

						localUser, createErr := userRepo.CreateFromClerk(r.Context(), clerkUserID, username, email)
						if createErr != nil {
							slog.Error("failed to create user on-demand",
								"clerk_user_id", clerkUserID,
								"error", createErr,
							)
							next.ServeHTTP(w, r)
							return
						}
						slog.Info("on-demand user creation",
							"clerk_user_id", clerkUserID,
							"local_user_id", localUser.ID,
						)
						user = localUser
					} else {
						slog.Error("failed to lookup user by clerk ID",
							"clerk_user_id", clerkUserID,
							"error", err,
						)
						next.ServeHTTP(w, r)
						return
					}
				}

				// Inject authenticated user into context
				authUser := &domain.AuthenticatedUser{
					ID:       user.ID,
					ClerkID:  user.ClerkUserID,
					Username: user.Username,
					Email:    user.Email,
					Role:     user.Role,
				}
				ctx := withUser(r.Context(), authUser)
				next.ServeHTTP(w, r.WithContext(ctx))
			}),
		)
	}
}
