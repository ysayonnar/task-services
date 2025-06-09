package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func New() Storage {
	return Storage{}
}

func (s *Storage) Conn() error {
	const op = "storage.Conn"

	dsn := os.Getenv("DSN")
	if dsn == "" {
		return fmt.Errorf("op: %s, err: %s", op, "no env dsn")
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	if err := conn.Ping(); err != nil {
		return fmt.Errorf("op: %s, ping error: %w", op, err)
	}
	s.DB = conn

	//TODO: переписать к хуям этот запрос на нужну схему
	query := `CREATE TABLE IF NOT EXISTS tasks(
		tasks_id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL;
		email VARCHAR(72) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = s.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("op: %s, table creating: %w", op, err)
	}

	return nil
}
