package database

import (
	"github.com/stretchr/testify/assert"
	"myapp/models"
	"testing"
)

// TestConnect tests the database connection
func TestConnect(t *testing.T) {
	// Call the Connect function
	Connect()

	// Assert that DB is not nil
	assert.NotNil(t, DB, "Database connection should not be nil")

	// Assert that the database is connected
	sqlDB, err := DB.DB()
	assert.NoError(t, err, "Getting database instance should not return an error")
	err = sqlDB.Ping()
	assert.NoError(t, err, "Database should be connected and pingable")
}

// TestMigrate tests the database migration
func TestMigrate(t *testing.T) {
	// Ensure database connection is established
	Connect()

	// Call the Migrate function
	Migrate()

	// Check if the tables exist in the database
	tables := []interface{}{
		&models.User{},
		&models.TopUp{},
		&models.Payment{},
		&models.Transfer{},
	}

	for _, table := range tables {
		assert.True(t, DB.Migrator().HasTable(table), "Table should exist in the database")
	}
}
