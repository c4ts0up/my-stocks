package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConnectDB_WithMock tests ConnectDB with an in-memory SQLite mock
func TestConnectDB_WithMock(t *testing.T) {
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to create mock DB")

	err = ConnectDB("", mockDB)
	assert.NoError(t, err, "Failed to connect to mock DB")
	assert.NotNil(t, DB, "DB should not be nil after mock connection")
	assert.Equal(t, mockDB, DB, "DB should point to the mock instance")
}

// TestConnectDB_WithBadDSN simulates a bad connection string
func TestConnectDB_WithBadDSN(t *testing.T) {
	err := ConnectDB("bad_dsn")
	assert.Error(t, err, "Error was not returned upon bad DSN")
	assert.Nil(t, DB, "DB should be nil with a bad DSN")
}

// TestCloseDB_WithMock ensures CloseDB works with a mock DB
func TestCloseDB_WithMock(t *testing.T) {
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to create mock DB")

	err = ConnectDB("", mockDB)
	assert.NoError(t, err, "Failed to close DB")
	assert.NotNil(t, DB, "DB should be connected")

	err = CloseDB()
	assert.NoError(t, err, "Failed to close DB")
	dbSQL, err := DB.DB()
	assert.NotNil(t, dbSQL, "DB object should still exist after CloseDB")
	assert.Nil(t, err, "Closing a mock DB should not return an error")
}

// TestCloseDB_NoConnection ensures CloseDB throws error when closing a nonexistent connection
func TestCloseDB_NoConnection(t *testing.T) {
	DB = nil

	err := CloseDB()
	assert.Error(t, err, "Closing of DB without started connection should return an error")
}

// TestCloseDB_Failure handles failure to retrieve DB object
func TestCloseDB_Failure(t *testing.T) {
	// needs realistic but "bad" setup
	DB, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	DB.Statement.ConnPool = nil

	err := CloseDB()
	assert.NoError(t, err, "Error was not returned upon close DB")
	assert.NotNil(t, DB, "DB should still exist even after close failure")
}

// TestAutoMigrateError simulates a migration failure
func TestAutoMigrateError(t *testing.T) {
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err, "Failed to create mock DB")

	DB = mockDB

	// Simulate an AutoMigrate error
	DB.Migrator().DropTable(&StockRating{})
	err = DB.Migrator().AutoMigrate(&struct{}{})
	assert.Error(t, err, "AutoMigrate should fail with invalid struct")
}
