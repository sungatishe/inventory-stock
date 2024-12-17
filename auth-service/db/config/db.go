package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var DB *sql.DB

func InitDb() {
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
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,                 -- Уникальный идентификатор
			email VARCHAR(255) NOT NULL UNIQUE,   -- Email пользователя
			password VARCHAR(255) NOT NULL,       -- Хэшированный пароль
			role VARCHAR(50) DEFAULT 'user',      -- Роль пользователя (по умолчанию 'user')
			created_at TIMESTAMP DEFAULT NOW(),   -- Дата создания записи
			updated_at TIMESTAMP DEFAULT NOW()    -- Дата последнего обновления
		);
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		panic(fmt.Sprintf("Failed to create users table: %v", err))
	}

}
