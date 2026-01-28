package resolvers

import (
	"github.com/yourorg/perspectize-go/internal/core/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	ContentService *services.ContentService
}

// NewResolver creates a new resolver with dependencies
func NewResolver(contentService *services.ContentService) *Resolver {
	return &Resolver{
		ContentService: contentService,
	}
}
