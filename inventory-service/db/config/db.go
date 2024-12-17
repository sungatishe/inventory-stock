package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func InitDB() {
	dsn := os.Getenv("dsn")
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("Failed to connect db ")
	}
	err = DB.Ping()

	if err != nil {
		panic("Failed to ping the database")
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		quantity INT NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		user_id VARCHAR(36) NOT NULL
	);
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		panic(fmt.Sprintf("Failed to create inventory table: %v", err))
	}
}
