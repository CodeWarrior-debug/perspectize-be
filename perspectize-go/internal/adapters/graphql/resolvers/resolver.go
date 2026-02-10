package resolvers

import (
	"github.com/CodeWarrior-debug/perspectize-be/perspectize-go/internal/core/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	ContentService     *services.ContentService
	UserService        *services.UserService
	PerspectiveService *services.PerspectiveService
}

// NewResolver creates a new resolver with dependencies
func NewResolver(
	contentService *services.ContentService,
	userService *services.UserService,
	perspectiveService *services.PerspectiveService,
) *Resolver {
	return &Resolver{
		ContentService:     contentService,
		UserService:        userService,
		PerspectiveService: perspectiveService,
	}
}
