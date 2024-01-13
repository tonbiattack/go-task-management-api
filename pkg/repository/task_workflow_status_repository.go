package repository

import (
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"gorm.io/gorm"
)

type TaskWorkflowStatusRepository struct {
	db *gorm.DB
}

func NewTaskWorkflowStatusRepository(db *gorm.DB) *TaskWorkflowStatusRepository {
	return &TaskWorkflowStatusRepository{db: db}
}

func (r *TaskWorkflowStatusRepository) Create(taskWorkflowStatus *model.TaskWorkflowStatus) error {
	return r.db.Create(taskWorkflowStatus).Error
}

func (r *TaskWorkflowStatusRepository) FindByID(id string) (*model.TaskWorkflowStatus, error) {
	var taskWorkflowStatus model.TaskWorkflowStatus
	err := r.db.Where("id = ?", id).First(&taskWorkflowStatus).Error
	return &taskWorkflowStatus, err
}

func (r *TaskWorkflowStatusRepository) Update(taskWorkflowStatus *model.TaskWorkflowStatus) error {
	return r.db.Save(taskWorkflowStatus).Error
}

func (r *TaskWorkflowStatusRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.TaskWorkflowStatus{}).Error
}

func (r *TaskWorkflowStatusRepository) FindByTaskID(taskID string) (*model.TaskWorkflowStatus, error) {
	var status model.TaskWorkflowStatus
	err := r.db.Where("task_id = ?", taskID).First(&status).Error
	return &status, err
}
