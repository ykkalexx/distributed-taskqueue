package task

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteQueue struct {
	db *sql.DB
}

func NewSQLiteQueue(dbPath string) (*SQLiteQueue, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create tasks table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			function_name TEXT,
			priority INTEGER,
			retries INTEGER,
			max_retries INTEGER,
			data TEXT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks table: %v", err)
	}

	return &SQLiteQueue{db: db}, nil
}

func (sq *SQLiteQueue) AddTask(task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}

	_, err = sq.db.Exec(
		"INSERT INTO tasks (id, function_name, priority, retries, max_retries, data) VALUES (?, ?, ?, ?, ?, ?)",
		task.ID, task.FunctionName, task.Priority, task.Retries, task.MaxRetries, string(data),
	)
	if err != nil {
		return fmt.Errorf("failed to insert task: %v", err)
	}

	return nil
}

func (sq *SQLiteQueue) GetTask() (Task, bool, error) {
	var task Task
	var data string

	err := sq.db.QueryRow(`
		SELECT data FROM tasks 
		ORDER BY priority DESC, retries ASC 
		LIMIT 1
	`).Scan(&data)

	if err == sql.ErrNoRows {
		return Task{}, false, nil
	} else if err != nil {
		return Task{}, false, fmt.Errorf("failed to get task: %v", err)
	}

	err = json.Unmarshal([]byte(data), &task)
	if err != nil {
		return Task{}, false, fmt.Errorf("failed to unmarshal task: %v", err)
	}

	_, err = sq.db.Exec("DELETE FROM tasks WHERE id = ?", task.ID)
	if err != nil {
		log.Printf("Failed to delete task %d from SQLite: %v", task.ID, err)
	}

	return task, true, nil
}

func (sq *SQLiteQueue) Close() error {
	return sq.db.Close()
}

// Ensure SQLiteQueue implements Queue interface
var _ Queue = (*SQLiteQueue)(nil)
