package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Создание таблиц
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS trainings (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description TEXT,
            session_id VARCHAR(50) NOT NULL
        );
        CREATE TABLE IF NOT EXISTS challenges (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS metrics (
            id SERIAL PRIMARY KEY,
            session_id VARCHAR(50) NOT NULL,
            pulse INT,
            speed FLOAT,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		return nil, err
	}

	return db, nil
}
