package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	host := os.Getenv("DB_POSTGRESQL_HOST")
	user := os.Getenv("DB_POSTGRESQL_USER")
	password := os.Getenv("DB_POSTGRESQL_PASSWORD")
	dbname := os.Getenv("DB_POSTGRESQL_NAME")
	port := os.Getenv("DB_POSTGRESQL_PORT")
	sslmode := os.Getenv("DB_POSTGRESQL_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to the database")
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Failed to get generic DB object: ")
		panic(err)
	}

	// Connection pool configuration
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Successfully connected to PostgreSQL with GORM")

	return db
}
