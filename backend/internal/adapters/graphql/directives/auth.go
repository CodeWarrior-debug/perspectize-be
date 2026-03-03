package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/web/middleware"
)

// DirectiveRoot holds the directive implementations for the GraphQL schema.
type DirectiveRoot struct{}

// NewDirectiveRoot creates a new DirectiveRoot with all directive handlers.
func NewDirectiveRoot() *DirectiveRoot {
	return &DirectiveRoot{}
}

// Auth enforces that the request is authenticated.
// Returns an error if no valid user context is present.
func (d *DirectiveRoot) Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	_, authenticated := middleware.ForContext(ctx)
	if !authenticated {
		return nil, fmt.Errorf("access denied: authentication required")
	}
	return next(ctx)
}

// Owner enforces that the request is authenticated and the user owns the resource.
// For now, only validates authentication — full ownership check deferred to Plan 02.
func (d *DirectiveRoot) Owner(ctx context.Context, obj interface{}, next graphql.Resolver, idField string) (interface{}, error) {
	_, authenticated := middleware.ForContext(ctx)
	if !authenticated {
		return nil, fmt.Errorf("access denied: authentication required")
	}

	// *TEMP* - Full ownership check (compare userID with resource owner) deferred to Plan 02
	return next(ctx)
}
