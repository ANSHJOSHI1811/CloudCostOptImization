package database

import (
	"fmt"
	"log"
	"cco_api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Database connection string
var DbConnStr = "host=localhost user=postgres password=12345678 dbname=temp_db sslmode=disable"

// Initialize Database
func InitDatabase() {
	var err error
	DB, err = gorm.Open(postgres.Open(DbConnStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&models.Provider{}, &models.Region{}, &models.SKU{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connected and migrated successfully.")
}