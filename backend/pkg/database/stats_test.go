package database_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CodeWarrior-debug/perspectize/backend/pkg/database"
)

func TestStatsHandler(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://invalid:invalid@localhost:5432/fake?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	handler := database.StatsHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/debug/db-stats", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var stats sql.DBStats
	err = json.NewDecoder(rec.Body).Decode(&stats)
	assert.NoError(t, err)
}
