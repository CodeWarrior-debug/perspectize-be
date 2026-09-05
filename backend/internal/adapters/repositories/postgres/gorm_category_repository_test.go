package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CodeWarrior-debug/perspectize/backend/internal/core/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var catRepoTime = time.Date(2026, 5, 6, 7, 8, 9, 0, time.UTC)

func categoryRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "wikidata_qid", "label", "description", "entity_type", "created_at", "updated_at",
	})
}

func TestGormCategoryRepository_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("returns mapped category", func(t *testing.T) {
		db, mock := newMockDB(t)
		mock.ExpectQuery(`SELECT \* FROM "categories"`).
			WillReturnRows(categoryRows().AddRow(3, "Q42", "Douglas Adams", "English author", "human", catRepoTime, catRepoTime))

		got, err := NewGormCategoryRepository(db).GetByID(ctx, 3)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, &domain.Category{
			ID: 3, WikidataQID: "Q42", Label: "Douglas Adams",
			Description: "English author", EntityType: "human",
			CreatedAt: catRepoTime, UpdatedAt: catRepoTime,
		}, got)
		assertAllExpectationsMet(t, mock)
	})

	t.Run("maps empty result to domain.ErrNotFound", func(t *testing.T) {
		db, mock := newMockDB(t)
		mock.ExpectQuery(`SELECT \* FROM "categories"`).WillReturnRows(categoryRows())

		got, err := NewGormCategoryRepository(db).GetByID(ctx, 404)
		assert.Nil(t, got)
		assert.True(t, errors.Is(err, domain.ErrNotFound), "expected domain.ErrNotFound, got %v", err)
		assertAllExpectationsMet(t, mock)
	})

	t.Run("wraps other driver errors with context", func(t *testing.T) {
		db, mock := newMockDB(t)
		mock.ExpectQuery(`SELECT \* FROM "categories"`).WillReturnError(errors.New("cat boom"))

		got, err := NewGormCategoryRepository(db).GetByID(ctx, 3)
		assert.Nil(t, got)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get category by id")
		assert.Contains(t, err.Error(), "cat boom")
		assertAllExpectationsMet(t, mock)
	})
}

func TestGormCategoryRepository_Upsert(t *testing.T) {
	ctx := context.Background()
	input := &domain.Category{WikidataQID: "Q42", Label: "Douglas Adams", Description: "English author", EntityType: "human"}

	t.Run("upserts on wikidata_qid conflict then re-reads the fresh row", func(t *testing.T) {
		db, mock := newMockDB(t)
		mock.ExpectQuery(`INSERT INTO "categories" .* ON CONFLICT`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
		mock.ExpectQuery(`SELECT \* FROM "categories" WHERE wikidata_qid`).
			WillReturnRows(categoryRows().AddRow(3, "Q42", "Douglas Adams", "English author", "human", catRepoTime, catRepoTime))

		got, err := NewGormCategoryRepository(db).Upsert(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, 3, got.ID)
		assert.Equal(t, "Q42", got.WikidataQID)
		assert.Equal(t, catRepoTime, got.UpdatedAt)
		assertAllExpectationsMet(t, mock)
	})

	t.Run("wraps upsert errors", func(t *testing.T) {
		db, mock := newMockDB(t)
		mock.ExpectQuery(`INSERT INTO "categories"`).WillReturnError(errors.New("upsert boom"))

		got, err := NewGormCategoryRepository(db).Upsert(ctx, input)
		assert.Nil(t, got)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to upsert category")
		assert.Contains(t, err.Error(), "upsert boom")
		assertAllExpectationsMet(t, mock)
	})

	t.Run("wraps re-read errors distinctly", func(t *testing.T) {
		db, mock := newMockDB(t)
		mock.ExpectQuery(`INSERT INTO "categories"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))
		mock.ExpectQuery(`SELECT \* FROM "categories"`).WillReturnError(errors.New("refetch boom"))

		got, err := NewGormCategoryRepository(db).Upsert(ctx, input)
		assert.Nil(t, got)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch upserted category")
		assert.Contains(t, err.Error(), "refetch boom")
		assertAllExpectationsMet(t, mock)
	})
}
