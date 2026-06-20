package database

import (
	"fmt"
	"jwt-auth-backend/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// NamingStrategy: schema.NamingStrategy{
		// 	SingularTable: true, // This disables pluralization DB Like add additional 's' Exlample: users, Devices
		// },

		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Device{},
		&models.Product{},
		&models.Customer{},
		&models.ContactPerson{},
		&models.ContactPersonInput{},
		&models.CustomerResponse{},
		&models.ContactPersonOutput{},
	)
	DB = db
	log.Println("Database connected & migrated")
}
