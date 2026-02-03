package domain

// ContentSortBy represents sortable fields for content queries
type ContentSortBy string

const (
	ContentSortByCreatedAt ContentSortBy = "CREATED_AT"
	ContentSortByUpdatedAt ContentSortBy = "UPDATED_AT"
	ContentSortByName      ContentSortBy = "NAME"
)

// SortOrder represents ascending or descending sort direction
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// ContentFilter contains filter criteria for content queries
type ContentFilter struct {
	ContentType      *ContentType
	MinLengthSeconds *int
	MaxLengthSeconds *int
}

// ContentListParams contains parameters for paginated content queries
type ContentListParams struct {
	First             *int
	After             *string // Opaque cursor (base64-encoded id)
	Last              *int
	Before            *string
	SortBy            ContentSortBy
	SortOrder         SortOrder
	IncludeTotalCount bool
	Filter            *ContentFilter
}

// PaginatedContent represents a paginated list of content
type PaginatedContent struct {
	Items       []*Content
	HasNext     bool
	HasPrev     bool
	StartCursor *string
	EndCursor   *string
	TotalCount  *int
}
