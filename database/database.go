package database

import (
	"database/sql"
	"fmt"
	"os"
)
import _ "github.com/lib/pq"

func Connect() (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println(dataSource)

	return sql.Open("postgres", dataSource)
}

func ConnectTest(connectionStr string) (*sql.DB, error) {
	return sql.Open("postgres", connectionStr)
}
