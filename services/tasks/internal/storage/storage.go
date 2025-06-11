package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"tasks/internal/models"

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

	query := `CREATE TABLE IF NOT EXISTS tasks(
		task_id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		deadline TIMESTAMP NOT NULL,
		is_notificate BOOL DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS categories(
		category_id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS tasks_categories(
		task_id INTEGER REFERENCES tasks(task_id) ON DELETE CASCADE,
		category_id INTEGER REFERENCES categories(category_id) ON DELETE CASCADE,
		PRIMARY KEY (task_id, category_id)
	);`

	_, err = s.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("op: %s, table creating: %w", op, err)
	}

	return nil
}

func (s *Storage) InsertTask(ctx context.Context, task models.Task)
