package services

import (
	"context"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
)

// WikidataClient defines the contract for Wikidata API interactions.
// Only wbsearchentities is supported — no SPARQL, no REST entity fetch.
type WikidataClient interface {
	// Search queries the Wikidata Entity Search API (wbsearchentities)
	Search(ctx context.Context, query string, language string, limit int) ([]domain.WikidataSearchResult, error)
}
