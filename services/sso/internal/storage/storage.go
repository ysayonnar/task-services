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
	defer conn.Close()

	s.DB = conn
	return nil
}
