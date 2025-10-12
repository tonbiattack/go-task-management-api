package repository

import (
	"strings"

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

// GetTaskWorkflowDetail retrieves a task with the current workflow step information using joins.
func (r *TaskRepository) GetTaskWorkflowDetail(id string) (*model.TaskWorkflowDetail, error) {
	var detail model.TaskWorkflowDetail
	columns := []string{
		"tasks.id AS task_id",
		"tasks.title AS task_title",
		"tasks.description AS task_description",
		"tasks.status AS task_status",
		"tasks.created_at AS task_created_at",
		"tasks.updated_at AS task_updated_at",
		"tasks.deadline AS task_deadline",
		"task_workflow_statuses.status AS workflow_status",
		"workflow_steps.id AS workflow_step_id",
		"workflow_steps.name AS workflow_step_name",
		"workflow_steps.description AS workflow_step_description",
		"workflow_steps.`order` AS workflow_step_order",
		"workflows.id AS workflow_id",
		"workflows.name AS workflow_name",
	}
	err := r.DB.Table("tasks").
		Select(strings.Join(columns, ", ")).
		Joins("JOIN task_workflow_statuses ON task_workflow_statuses.task_id = tasks.id").
		Joins("JOIN workflow_steps ON workflow_steps.id = task_workflow_statuses.workflow_step_id").
		Joins("JOIN workflows ON workflows.id = workflow_steps.workflow_id").
		Where("tasks.id = ?", id).
		Take(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}
