package repository

import (
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// CreateTask は新しいタスクをデータベースに挿入します。
func (r *TaskRepository) CreateTask(task *model.Task) error {
	return r.DB.Create(task).Error
}

// GetTaskByID はIDによってタスクをデータベースから取得します。
// 指定されたIDのタスクが存在しない場合、エラーを返します。
func (r *TaskRepository) GetTaskByID(id string) (*model.Task, error) {
	var task model.Task
	err := r.DB.First(&task, "id = ?", id).Error
	return &task, err
}

// GetAllTasks はデータベース内のすべてのタスクを取得します。
func (r *TaskRepository) GetAllTasks() ([]*model.Task, error) {
	var tasks []*model.Task
	err := r.DB.Find(&tasks).Error
	return tasks, err
}

// UpdateTask は指定されたIDのタスクを更新します。
// タスクが存在しない場合、エラーを返します。
func (r *TaskRepository) UpdateTask(task *model.Task) error {
	return r.DB.Save(task).Error
}

// DeleteTask は指定されたIDのタスクをデータベースから削除します。
// タスクが存在しない場合、エラーを返します。
func (r *TaskRepository) DeleteTask(id string) error {
	return r.DB.Delete(&model.Task{}, id).Error
}
