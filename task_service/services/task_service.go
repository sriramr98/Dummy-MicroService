package services

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type TaskService struct {
	db *sql.DB
}

type Task struct {
	Id    int64
	Title string
	Body  string
}

func NewTaskService(db *sql.DB) TaskService {
	return TaskService{db: db}
}

func (ts TaskService) CreateTask(title, body string, userId int64) (int64, error) {
	var id int64
	err := ts.db.QueryRow("INSERT INTO tasks (title, body, userId) VALUES ($1, $2, $3) RETURNING id", title, body, userId).Scan(&id)
	return id, err
}

func (ts TaskService) UpdateTask(id int64, title, body string, userId int64) error {
	_, err := ts.db.Exec("UPDATE tasks SET title = $1, body = $2 WHERE id = $3 AND userId = $4", title, body, id, userId)
	return err
}

func (ts TaskService) GetTask(id int64, userId int64) (Task, error) {
	var task Task
	err := ts.db.QueryRow("SELECT id, title, body FROM tasks WHERE id = $1 AND userId = $2", id, userId).Scan(&task.Id, &task.Title, &task.Body)
	return task, err
}

func (ts TaskService) ListTasks(userId int64) ([]Task, error) {
	rows, err := ts.db.Query("SELECT id, title, body FROM tasks WHERE userId = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Title, &task.Body); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (ts TaskService) DeleteTask(id int, userId int64) error {
	_, err := ts.db.Exec("DELETE FROM tasks WHERE id = $1 AND userId = $2", id, userId)
	return err
}
