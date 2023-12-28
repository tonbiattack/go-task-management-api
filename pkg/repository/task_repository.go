package repository

import (
	"database/sql"

	"github.com/tonbiattack/go-task-management-api/pkg/model"
)

// TaskRepository はタスクストレージへのアクセスを提供します。
type TaskRepository struct {
	DB *sql.DB
}

// NewTaskRepository は新しいTaskRepositoryを与えられたデータベース接続で作成します。
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// CreateTask は新しいタスクをデータベースに挿入します。
// タスクステータスが設定されていない場合、空の状態で省略されることを意味します。
func (r *TaskRepository) CreateTask(task *model.Task) error {
	var query string
	var err error

	if task.Status == "" {
		// ステータスを省略
		query = `INSERT INTO Task (id, title, description, deadline) VALUES (?, ?, ?, ?)`
		_, err = r.DB.Exec(query, task.ID, task.Title, task.Description, task.Deadline)
	} else {
		// ステータスを含める
		query = `INSERT INTO Task (id, title, description, status, deadline) VALUES (?, ?, ?, ?, ?)`
		_, err = r.DB.Exec(query, task.ID, task.Title, task.Description, task.Status, task.Deadline)
	}
	return err
}

// GetTaskByID はIDによってタスクをデータベースから取得します。
// 指定されたIDのタスクが存在しない場合、エラーを返します。
func (r *TaskRepository) GetTaskByID(id string) (*model.Task, error) {
	var task model.Task
	query := `SELECT id, title, description, status, created_at, updated_at, deadline FROM Task WHERE id = ?`
	err := r.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.Deadline)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllTasks はデータベース内のすべてのタスクを取得します。
func (r *TaskRepository) GetAllTasks() ([]*model.Task, error) {
	tasks := []*model.Task{}

	query := `SELECT id, title, description, status, created_at, updated_at, deadline FROM Task`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.Deadline); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// UpdateTask は指定されたIDのタスクを更新します。
// タスクが存在しない場合、エラーを返します。
func (r *TaskRepository) UpdateTask(task *model.Task) error {
	query := `UPDATE Task SET title = ?, description = ?, status = ?, updated_at = NOW(), deadline = ? WHERE id = ?`
	res, err := r.DB.Exec(query, task.Title, task.Description, task.Status, task.Deadline, task.ID)
	if err != nil {
		return err
	}

	// 更新された行の数を確認
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteTask は指定されたIDのタスクをデータベースから削除します。
// タスクが存在しない場合、エラーを返します。
func (r *TaskRepository) DeleteTask(id string) error {
	query := `DELETE FROM Task WHERE id = ?`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	// 削除された行の数を確認
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
