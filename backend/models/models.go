package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=root password=mysecret dbname=stocks port=26257 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Connected to the database!")

	// Auto-migrate schemas
	err = DB.AutoMigrate(&StockRating{})
	if err != nil {
		return
	}
}

func CloseDB() {
	dbSQL, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve DB object: %v", err)
	}
	err = dbSQL.Close()
	if err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
		return
	}
}
