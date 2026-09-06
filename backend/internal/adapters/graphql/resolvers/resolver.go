package resolvers

import (
	"github.com/CodeWarrior-debug/perspectize/backend/internal/adapters/realtime"
	portservices "github.com/CodeWarrior-debug/perspectize/backend/internal/core/ports/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	ContentService     portservices.ContentService
	UserService        portservices.UserService
	PerspectiveService portservices.PerspectiveService
	CategoryService    portservices.CategoryService

	// Messaging feature dependencies.
	Messaging portservices.MessagingService
	Hub       *realtime.Hub
	Presence  *realtime.PresenceTracker
}

// NewResolver creates a new resolver with dependencies
func NewResolver(
	contentService portservices.ContentService,
	userService portservices.UserService,
	perspectiveService portservices.PerspectiveService,
	categoryService portservices.CategoryService,
	messagingService portservices.MessagingService,
	hub *realtime.Hub,
	presence *realtime.PresenceTracker,
) *Resolver {
	return &Resolver{
		ContentService:     contentService,
		UserService:        userService,
		PerspectiveService: perspectiveService,
		CategoryService:    categoryService,
		Messaging:          messagingService,
		Hub:                hub,
		Presence:           presence,
	}
}
