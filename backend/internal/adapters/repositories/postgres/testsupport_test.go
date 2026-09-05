package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newMockDB returns a *gorm.DB backed by a go-sqlmock *sql.DB.
//
// Three gorm.Config settings are load-bearing and must not be changed:
//   - DisableAutomaticPing: gorm.Open pings ConnPool at open time; sqlmock would
//     report "call to Ping was not expected".
//   - SkipDefaultTransaction: without it every Create/Update/Delete is wrapped in
//     BEGIN/COMMIT and each test would need ExpectBegin/ExpectCommit.
//   - Logger silent: GORM otherwise prints every statement plus every
//     ErrRecordNotFound as a warning, which is expected in these tests.
//
// gorm.io/driver/postgres performs no network round-trip during Initialize when
// Config.Conn is set, so no real database is required.
//
// QueryMatcherRegexp is used because GORM emits quoted identifiers
// (`SELECT * FROM "users" WHERE ...`); the default matcher requires exact
// whitespace-normalised equality, which is far too brittle here. Expectation
// patterns are unanchored Go regexps, so remember to escape `(`, `)`, `*`, `?`
// and `.` — or match on a distinctive substring instead.
func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	gdb, err := gorm.Open(
		gormpg.New(gormpg.Config{Conn: sqlDB}),
		&gorm.Config{
			DisableAutomaticPing:   true,
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Silent),
		},
	)
	require.NoError(t, err)

	return gdb, mock
}

// assertAllExpectationsMet fails the test if any queued sqlmock expectation was
// never consumed — this is what proves the repository actually issued the query.
func assertAllExpectationsMet(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewMockDB_OpensWithoutARealDatabase(t *testing.T) {
	db, mock := newMockDB(t)
	require.NotNil(t, db)
	require.NotNil(t, mock)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "role", "active"}).
			AddRow(1, "alice", "admin", true))

	var got UserModel
	require.NoError(t, db.Where("id = ?", 1).First(&got).Error)

	assert.Equal(t, 1, got.ID)
	assert.Equal(t, "alice", got.Username)
	assert.Equal(t, "admin", got.Role)
	assert.True(t, got.Active)
	assertAllExpectationsMet(t, mock)
}
