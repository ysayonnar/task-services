package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sso/internal/models"

	_ "github.com/lib/pq"
)

var ErrUserExists = errors.New("user already exists")

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
	return nil
}

// TODO: implement
func (s *Storage) GetUserByEmail(email string) (models.User, error) {
	return models.User{}, nil
}

func (s *Storage) InsertUser(email string, passwordHash string) (int64, error) {

}

func (s *Storage) DeleteUser(email string) (int64, error) {
	return 0, nil
}
