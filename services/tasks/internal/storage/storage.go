package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	tasks "github.com/ysayonnar/task-contracts/tasks/gen/go"
	"os"
	"tasks/internal/models"
)

var ErrCategoryNotFound = errors.New("category with such id doesn't exist")
var ErrCategoryAlreadyExists = errors.New("category with such name already exists")
var ErrTaskNotFound = errors.New("task with such id was not found")

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
		name VARCHAR(255) UNIQUE NOT NULL
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

func (s *Storage) InsertTask(ctx context.Context, task models.Task, categoryId int64) (int64, error) {
	const op = "storage.InsertTask"

	tx, err := s.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	var taskId int64
	queryTask := `INSERT INTO tasks(user_id, title, description, deadline, is_notificate) VALUES($1, $2, $3, $4, $5) RETURNING task_id;`

	err = tx.QueryRowContext(ctx, queryTask, task.UserId, task.Title, task.Description, task.Deadline, task.IsNotificate).Scan(&taskId)
	if err != nil {
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	queryLink := `INSERT INTO tasks_categories(task_id, category_id) VALUES($1, $2);`
	_, err = tx.ExecContext(ctx, queryLink, taskId, categoryId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return 0, ErrCategoryNotFound
			}
		}
		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return taskId, nil
}

func (s *Storage) InsertCategory(ctx context.Context, name string) (int64, error) {
	const op = "storage.InsertCategory"

	query := `INSERT INTO categories(name) VALUES($1) RETURNING category_id;`

	var categoryId int64
	err := s.DB.QueryRowContext(ctx, query, name).Scan(&categoryId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return 0, ErrCategoryAlreadyExists
			}
		}

		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return categoryId, nil
}

func (s *Storage) DeleteTask(ctx context.Context, userId int64, taskId int64) (int64, error) {
	const op = "storage.DeleteTask"

	query := `DELETE FROM tasks WHERE user_id = $1 AND task_id = $2 RETURNING task_id;`

	var deletedTaskId int64
	err := s.DB.QueryRowContext(ctx, query, userId, taskId).Scan(&deletedTaskId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrTaskNotFound
		}

		return 0, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return deletedTaskId, nil
}

func (s *Storage) GetTasksByUserId(ctx context.Context, userId int64) ([]*tasks.Task, error) {
	const op = "storage.GetTasksByUserId"

	query := `SELECT (t.task_id, tc.category_id, c.name, t.title, t.description, t.deadline, t.is_notificate, t.created_at) FROM tasks AS t INNER JOIN tasks_categories AS tc ON t.task_id = tc.task_id INNER JOIN categories AS c ON c.category_id = tc.category_id WHERE t.task_id = $1;`

	rows, err := s.DB.QueryContext(ctx, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}
	defer rows.Close()

	var selectedTasks []*tasks.Task
	for rows.Next() {
		var task tasks.Task

		err := rows.Scan(
			&task.TaskId,
			&task.CategoryId,
			&task.CategoryName,
			&task.Title,
			&task.Description,
			&task.Deadline,
			&task.IsNotificate,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("op: %s, err: %w", op, err)
		}
		selectedTasks = append(selectedTasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("op: %s, err: %w", op, err)
	}

	return selectedTasks, nil
}
