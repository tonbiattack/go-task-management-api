package model

import "time"

// TaskWorkflowDetail aggregates task information with the associated workflow step.
type TaskWorkflowDetail struct {
	TaskID            string     `json:"task_id" gorm:"column:task_id"`
	TaskTitle         string     `json:"task_title" gorm:"column:task_title"`
	TaskDescription   string     `json:"task_description" gorm:"column:task_description"`
	TaskStatus        *string    `json:"task_status" gorm:"column:task_status"`
	TaskCreatedAt     time.Time  `json:"task_created_at" gorm:"column:task_created_at"`
	TaskUpdatedAt     time.Time  `json:"task_updated_at" gorm:"column:task_updated_at"`
	TaskDeadline      *time.Time `json:"task_deadline" gorm:"column:task_deadline"`
	WorkflowStatus    string     `json:"workflow_status" gorm:"column:workflow_status"`
	WorkflowID        string     `json:"workflow_id" gorm:"column:workflow_id"`
	WorkflowName      string     `json:"workflow_name" gorm:"column:workflow_name"`
	WorkflowStepID    string     `json:"workflow_step_id" gorm:"column:workflow_step_id"`
	WorkflowStepName  string     `json:"workflow_step_name" gorm:"column:workflow_step_name"`
	WorkflowStepOrder int        `json:"workflow_step_order" gorm:"column:workflow_step_order"`
	WorkflowStepDesc  string     `json:"workflow_step_description" gorm:"column:workflow_step_description"`
}
