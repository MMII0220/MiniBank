package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sqlx.DB

// Connection to Database
func InitDB() (*sqlx.DB, error) {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("no env file found: %v", err)
	}

	dbc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = sqlx.Connect("postgres", dbc)
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	return db, nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
