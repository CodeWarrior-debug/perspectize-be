package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourorg/perspectize-go/pkg/database"
)

func TestConnect_ValidDSN(t *testing.T) {
	// Use test database connection
	dsn := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"

	db, err := database.Connect(dsn)

	if err != nil {
		t.Skip("Skipping test - PostgreSQL not available. Run 'make docker-up' to start database.")
	}

	assert.NoError(t, err)
	assert.NotNil(t, db)

	defer db.Close()

	// Verify connection works
	err = database.Ping(context.Background(), db)
	assert.NoError(t, err)
}

func TestConnect_InvalidDSN(t *testing.T) {
	dsn := "host=invalid port=9999 user=fake password=fake dbname=fake sslmode=disable"

	db, err := database.Connect(dsn)

	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestPing_Success(t *testing.T) {
	dsn := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"

	db, err := database.Connect(dsn)

	if err != nil {
		t.Skip("Skipping test - PostgreSQL not available")
	}

	defer db.Close()

	ctx := context.Background()
	err = database.Ping(ctx, db)

	assert.NoError(t, err)
}
