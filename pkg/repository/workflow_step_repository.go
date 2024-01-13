package repository

import (
	"errors"

	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"gorm.io/gorm"
)

type WorkflowStepRepository struct {
	db *gorm.DB
}

func NewWorkflowStepRepository(db *gorm.DB) *WorkflowStepRepository {
	return &WorkflowStepRepository{db: db}
}

func (r *WorkflowStepRepository) Create(workflowStep *model.WorkflowStep) error {
	return r.db.Create(workflowStep).Error
}

func (r *WorkflowStepRepository) FindByID(id string) (*model.WorkflowStep, error) {
	var workflowStep model.WorkflowStep
	err := r.db.Where("id = ?", id).First(&workflowStep).Error
	return &workflowStep, err
}

func (r *WorkflowStepRepository) Update(workflowStep *model.WorkflowStep) error {
	return r.db.Save(workflowStep).Error
}

func (r *WorkflowStepRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.WorkflowStep{}).Error
}

// GetNextStep は、指定されたワークフローステップIDに基づいて次のステップを取得します。
func (r *WorkflowStepRepository) GetNextStep(currentStepID string) (*model.WorkflowStep, error) {
	var currentStep, nextStep model.WorkflowStep

	// 現在のステップを取得
	if err := r.db.Where("id = ?", currentStepID).First(&currentStep).Error; err != nil {
		return nil, err
	}

	// 次のステップを取得（現在のステップの order より大きい最小の order を持つステップ）

	err := r.db.Where("workflow_id = ? AND `order` > ?", currentStep.WorkflowID, currentStep.Order).Order("`order` ASC").First(&nextStep).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoNextStep // ErrNoNextStep はカスタムエラーです。
	}
	return &nextStep, err
}

// GetInitialStep - 特定のワークフローIDに対する最初のステップを取得
func (r *WorkflowStepRepository) GetInitialStep(workflowID string) (*model.WorkflowStep, error) {
	var step model.WorkflowStep
	if err := r.db.Where("workflow_id = ?", workflowID).Order("`order` ASC").First(&step).Error; err != nil {
		return nil, err
	}
	return &step, nil
}
