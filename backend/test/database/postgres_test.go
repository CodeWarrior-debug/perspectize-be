package database_test

import (
	"context"
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestConnectGORM_ValidDSN(t *testing.T) {
	// Use test database connection
	dsn := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	poolCfg := database.DefaultPoolConfig()

	db, err := database.ConnectGORM(dsn, poolCfg)

	if err != nil {
		t.Skip("Skipping test - PostgreSQL not available. Run 'make docker-up' to start database.")
	}

	assert.NoError(t, err)
	assert.NotNil(t, db)

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Verify connection works
	err = database.PingGORM(context.Background(), db)
	assert.NoError(t, err)
}

func TestConnectGORM_InvalidDSN(t *testing.T) {
	dsn := "host=invalid port=9999 user=fake password=fake dbname=fake sslmode=disable"
	poolCfg := database.DefaultPoolConfig()

	db, err := database.ConnectGORM(dsn, poolCfg)

	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to initialize GORM")
}

func TestPingGORM_Success(t *testing.T) {
	dsn := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	poolCfg := database.DefaultPoolConfig()

	db, err := database.ConnectGORM(dsn, poolCfg)

	if err != nil {
		t.Skip("Skipping test - PostgreSQL not available")
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	ctx := context.Background()
	err = database.PingGORM(ctx, db)

	assert.NoError(t, err)
}
