package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"myapp/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Connected to the database successfully!")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.TopUp{},
		&models.Payment{},
		&models.Transfer{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate models: %v\n", err)
	}

	fmt.Println("Database migrated successfully!")
}
