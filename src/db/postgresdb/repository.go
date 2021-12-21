package postgresdb

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func PostgresInit() {
	// host, user, password, name, port := os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"), os.Getenv("DB_PORT")
	host, user, password, name, port := "localhost", "postgres", "password", "server_db", "5433"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
