package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// DB holds the global database connection
var DB *gorm.DB

// ConnectDB establishes a database connection (default or mockable)
func ConnectDB(dsn string, dbInstance ...*gorm.DB) error {
	var err error

	// Use a provided DB instance (e.g., a mock) if available
	if len(dbInstance) > 0 {
		DB = dbInstance[0]
		log.Println("Using given database instance")
	} else {
		log.Println("Creating new database instance")
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			DB = nil
			return fmt.Errorf("failed to connect to database: %v", err)
		}
	}

	log.Println("✅ Connected to the database!")

	// Auto-migrate schemas
	err = DB.AutoMigrate(&Stock{}, &StockRating{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %v", err)
	}

	return nil
}

// CloseDB closes the database connection (default or mockable)
func CloseDB() error {
	if DB == nil {
		return fmt.Errorf("there is no open database connection")
	}
	// DB.DB() could panic, but let's hope it doesn't happen as only this class handles it
	dbSQL, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve DB object: %v", err)
	}
	err = dbSQL.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %v", err)
	}

	return nil
}
