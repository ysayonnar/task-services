package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sso/internal/models"

	"github.com/lib/pq"
)

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user was not found")

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

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.GetUserByEmail"

	query := `SELECT * FROM users WHERE email = $1;`

	var user models.User
	err := s.DB.QueryRowContext(ctx, query, email).Scan(&user.UserId, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return user, nil
}

func (s *Storage) InsertUser(ctx context.Context, email string, passwordHash string) (int64, error) {
	const op = "storage.InsertUser"

	query := `INSERT INTO users(email, password_hash) VALUES($1, $2) RETURNING user_id;`

	var userId int64
	err := s.DB.QueryRowContext(ctx, query, email, passwordHash).Scan(&userId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return 0, ErrUserExists
			}
		}
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return userId, nil
}

func (s *Storage) DeleteUser(ctx context.Context, email string) (int64, error) {
	const op = "storage.DeleteUser"

	query := `DELETE FROM users WHERE email = $1 RETURNING user_id;`

	var userId int64
	err := s.DB.QueryRowContext(ctx, query, email).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return userId, nil
}
