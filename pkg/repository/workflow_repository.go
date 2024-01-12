package repository

import (
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"gorm.io/gorm"
)

type WorkflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) Create(workflow *model.Workflow) error {
	return r.db.Create(workflow).Error
}

func (r *WorkflowRepository) FindByID(id string) (*model.Workflow, error) {
	var workflow model.Workflow
	err := r.db.Where("id = ?", id).First(&workflow).Error
	return &workflow, err
}

func (r *WorkflowRepository) Update(workflow *model.Workflow) error {
	return r.db.Save(workflow).Error
}

func (r *WorkflowRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Workflow{}).Error
}
