package model

import "time"

type TaskWorkflowStatus struct {
	ID             string       `gorm:"type:char(36);primaryKey" json:"id"`
	TaskID         string       `gorm:"type:char(36);not null" json:"task_id"`
	WorkflowStepID string       `gorm:"type:char(36);not null" json:"workflow_step_id"`
	Status         string       `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Task           Task         `gorm:"foreignKey:TaskID"`
	WorkflowStep   WorkflowStep `gorm:"foreignKey:WorkflowStepID"`
}
