package database

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

const slowQueryThreshold = 100 * time.Millisecond

// RegisterSlowQueryLogger adds GORM callbacks that log queries exceeding the threshold.
func RegisterSlowQueryLogger(db *gorm.DB) {
	callback := func(operation string) func(db *gorm.DB) {
		return func(db *gorm.DB) {
			if db.Statement == nil {
				return
			}
			startVal := db.Statement.Context.Value(queryStartKey{})
			if startVal == nil {
				return
			}
			duration := time.Since(startVal.(time.Time))
			if duration >= slowQueryThreshold {
				slog.WarnContext(db.Statement.Context, "slow query",
					"operation", operation,
					"duration_ms", duration.Milliseconds(),
					"sql", db.Statement.SQL.String(),
					"rows", db.Statement.RowsAffected,
				)
			}
		}
	}

	startTimer := func(db *gorm.DB) {
		db.Statement.Context = context.WithValue(db.Statement.Context, queryStartKey{}, time.Now())
	}

	_ = db.Callback().Query().Before("gorm:query").Register("perf:query_start", startTimer)
	_ = db.Callback().Query().After("gorm:query").Register("perf:query_slow", callback("query"))
	_ = db.Callback().Create().Before("gorm:create").Register("perf:create_start", startTimer)
	_ = db.Callback().Create().After("gorm:create").Register("perf:create_slow", callback("create"))
	_ = db.Callback().Update().Before("gorm:update").Register("perf:update_start", startTimer)
	_ = db.Callback().Update().After("gorm:update").Register("perf:update_slow", callback("update"))
	_ = db.Callback().Delete().Before("gorm:delete").Register("perf:delete_start", startTimer)
	_ = db.Callback().Delete().After("gorm:delete").Register("perf:delete_slow", callback("delete"))
}

type queryStartKey struct{}
