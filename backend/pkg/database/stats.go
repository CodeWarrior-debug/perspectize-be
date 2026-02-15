package database

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// StatsHandler returns an http.HandlerFunc that reports sql.DBStats as JSON.
func StatsHandler(sqlDB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := sqlDB.Stats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
