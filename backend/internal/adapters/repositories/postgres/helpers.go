package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	paginator "github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

// JSONBArray is a custom type for PostgreSQL jsonb[] columns.
// It stores JSON data as strings and properly handles NULL values.
type JSONBArray []string

// Scan implements the sql.Scanner interface for reading jsonb[] from PostgreSQL
func (a *JSONBArray) Scan(src interface{}) error {
	if src == nil {
		*a = nil
		return nil
	}

	// PostgreSQL returns jsonb[] as []byte containing array of bytea
	// We need to scan through StringArray
	var strArray StringArray
	if err := strArray.Scan(src); err != nil {
		return err
	}
	*a = JSONBArray(strArray)
	return nil
}

// Value implements the driver.Valuer interface for writing jsonb[] to PostgreSQL
func (a JSONBArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return StringArray(a).Value()
}

// encodeCursor encodes a content ID as an opaque base64 cursor
func encodeCursor(id int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor:%d", id)))
}

// decodeCursor decodes an opaque base64 cursor back to a content ID
func decodeCursor(cursor string) (int, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor: %w", err)
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 || parts[0] != "cursor" {
		return 0, fmt.Errorf("invalid cursor format")
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid cursor id: %w", err)
	}
	return id, nil
}

// sortColumnName maps a domain sort field to a safe SQL column name or JSONB extraction
func sortColumnName(sortBy domain.ContentSortBy) string {
	switch sortBy {
	case domain.ContentSortByUpdatedAt:
		return "updated_at"
	case domain.ContentSortByName:
		return "name"
	case domain.ContentSortByCreatedAt:
		return "created_at"
	case domain.ContentSortByViewCount:
		return "(response->'items'->0->'statistics'->>'viewCount')::BIGINT"
	case domain.ContentSortByLikeCount:
		return "(response->'items'->0->'statistics'->>'likeCount')::BIGINT"
	case domain.ContentSortByPublishedAt:
		return "response->'items'->0->'snippet'->>'publishedAt'"
	default:
		return "created_at"
	}
}

// sortDirection returns a safe SQL sort direction string
func sortDirection(order domain.SortOrder) string {
	if order == domain.SortOrderAsc {
		return "ASC"
	}
	return "DESC"
}

// isJSONBSort returns true if the sort field requires JSONB extraction
func isJSONBSort(sortBy domain.ContentSortBy) bool {
	return sortBy == domain.ContentSortByViewCount ||
		sortBy == domain.ContentSortByLikeCount ||
		sortBy == domain.ContentSortByPublishedAt
}

// perspectiveSortColumnName maps a domain sort field to a safe SQL column name
func perspectiveSortColumnName(sortBy domain.PerspectiveSortBy) string {
	switch sortBy {
	case domain.PerspectiveSortByUpdatedAt:
		return "updated_at"
	case domain.PerspectiveSortByClaim:
		return "claim"
	case domain.PerspectiveSortByCreatedAt:
		return "created_at"
	default:
		return "created_at"
	}
}

// contentTypeToDBValue converts domain ContentType to lowercase for database storage
func contentTypeToDBValue(ct domain.ContentType) string {
	return strings.ToLower(string(ct))
}

// contentTypeFromDBValue converts lowercase database value to domain ContentType
func contentTypeFromDBValue(s string) domain.ContentType {
	return domain.ContentType(strings.ToUpper(s))
}

// privacyToDBValue converts domain Privacy to lowercase for database storage
func privacyToDBValue(p domain.Privacy) string {
	return strings.ToLower(string(p))
}

// privacyFromDBValue converts lowercase database value to domain Privacy
func privacyFromDBValue(s string) domain.Privacy {
	return domain.Privacy(strings.ToUpper(s))
}

// reviewStatusToDBValue converts domain ReviewStatus to lowercase for database storage
func reviewStatusToDBValue(rs *domain.ReviewStatus) sql.NullString {
	if rs == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.ToLower(string(*rs)), Valid: true}
}

// reviewStatusFromDBValue converts lowercase database value to domain ReviewStatus
func reviewStatusFromDBValue(s sql.NullString) *domain.ReviewStatus {
	if !s.Valid {
		return nil
	}
	rs := domain.ReviewStatus(strings.ToUpper(s.String))
	return &rs
}

// intSliceToInt64Array converts []int to Int64Array for database storage
func intSliceToInt64Array(ints []int) Int64Array {
	if ints == nil {
		return nil
	}
	result := make(Int64Array, len(ints))
	for i, v := range ints {
		result[i] = int64(v)
	}
	return result
}

// buildContentSortRules builds paginator rules for content sorting
// Returns slice with primary sort rule + ID tie-breaker rule
func buildContentSortRules(sortBy domain.ContentSortBy, order domain.SortOrder) []paginator.Rule {
	// Map domain.SortOrder to paginator.Order
	var paginatorOrder paginator.Order
	if order == domain.SortOrderAsc {
		paginatorOrder = paginator.ASC
	} else {
		paginatorOrder = paginator.DESC
	}

	var primaryRule paginator.Rule

	switch sortBy {
	case domain.ContentSortByViewCount:
		primaryRule = paginator.Rule{
			Key:             "ViewCount",
			Order:           paginatorOrder,
			SQLRepr:         "(response->'items'->0->'statistics'->>'viewCount')::BIGINT",
			NULLReplacement: int64(0),
		}
	case domain.ContentSortByLikeCount:
		primaryRule = paginator.Rule{
			Key:             "LikeCount",
			Order:           paginatorOrder,
			SQLRepr:         "(response->'items'->0->'statistics'->>'likeCount')::BIGINT",
			NULLReplacement: int64(0),
		}
	case domain.ContentSortByPublishedAt:
		primaryRule = paginator.Rule{
			Key:             "PublishedAt",
			Order:           paginatorOrder,
			SQLRepr:         "response->'items'->0->'snippet'->>'publishedAt'",
			NULLReplacement: "",
		}
	case domain.ContentSortByUpdatedAt:
		primaryRule = paginator.Rule{
			Key:   "UpdatedAt",
			Order: paginatorOrder,
		}
	case domain.ContentSortByName:
		primaryRule = paginator.Rule{
			Key:   "Name",
			Order: paginatorOrder,
		}
	case domain.ContentSortByCreatedAt:
		primaryRule = paginator.Rule{
			Key:   "CreatedAt",
			Order: paginatorOrder,
		}
	default:
		// Default to CreatedAt DESC
		primaryRule = paginator.Rule{
			Key:   "CreatedAt",
			Order: paginator.DESC,
		}
	}

	// Tie-breaker: ID with same sort direction as primary
	tieBreaker := paginator.Rule{
		Key:   "ID",
		Order: paginatorOrder,
	}

	return []paginator.Rule{primaryRule, tieBreaker}
}

// buildPerspectiveSortRules builds paginator rules for perspective sorting
// Returns slice with primary sort rule + ID tie-breaker rule
func buildPerspectiveSortRules(sortBy domain.PerspectiveSortBy, order domain.SortOrder) []paginator.Rule {
	// Map domain.SortOrder to paginator.Order
	var paginatorOrder paginator.Order
	if order == domain.SortOrderAsc {
		paginatorOrder = paginator.ASC
	} else {
		paginatorOrder = paginator.DESC
	}

	var primaryRule paginator.Rule

	switch sortBy {
	case domain.PerspectiveSortByUpdatedAt:
		primaryRule = paginator.Rule{
			Key:   "UpdatedAt",
			Order: paginatorOrder,
		}
	case domain.PerspectiveSortByClaim:
		primaryRule = paginator.Rule{
			Key:   "Claim",
			Order: paginatorOrder,
		}
	case domain.PerspectiveSortByCreatedAt:
		primaryRule = paginator.Rule{
			Key:   "CreatedAt",
			Order: paginatorOrder,
		}
	default:
		// Default to CreatedAt DESC
		primaryRule = paginator.Rule{
			Key:   "CreatedAt",
			Order: paginator.DESC,
		}
	}

	// Tie-breaker: ID with same sort direction as primary
	tieBreaker := paginator.Rule{
		Key:   "ID",
		Order: paginatorOrder,
	}

	return []paginator.Rule{primaryRule, tieBreaker}
}
