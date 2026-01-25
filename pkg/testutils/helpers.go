package testutils

import (
	"testing"
	"time"

	"github.com/tonbiattack/go-task-management-api/pkg/config"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"gorm.io/gorm"
)

// SetupTestDB テスト用のデータベースをセットアップし、テスト終了時にクリーンアップします
func SetupTestDB(t *testing.T) *gorm.DB {
	db := config.InitTestDB()

	t.Cleanup(func() {
		// テスト終了時のクリーンアップ
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})

	return db
}

// CreateTestTask テスト用のタスクを作成します
func CreateTestTask(db *gorm.DB, id, title, description string) *model.Task {
	status := "pending"
	task := &model.Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      &status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(task)
	return task
}

// CreateTestWorkflow テスト用のワークフローを作成します
func CreateTestWorkflow(db *gorm.DB, id, name, description string) *model.Workflow {
	workflow := &model.Workflow{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(workflow)
	return workflow
}

// CreateTestWorkflowStep テスト用のワークフローステップを作成します
func CreateTestWorkflowStep(db *gorm.DB, id, workflowID, name string, order int) *model.WorkflowStep {
	step := &model.WorkflowStep{
		ID:          id,
		WorkflowID:  workflowID,
		Name:        name,
		Description: "Test step",
		Order:       order,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.Create(step)
	return step
}

// CreateTestTaskWorkflowStatus テスト用のタスクワークフローステータスを作成します
func CreateTestTaskWorkflowStatus(db *gorm.DB, id, taskID, workflowStepID, status string) *model.TaskWorkflowStatus {
	taskStatus := &model.TaskWorkflowStatus{
		ID:             id,
		TaskID:         taskID,
		WorkflowStepID: workflowStepID,
		Status:         status,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(taskStatus)
	return taskStatus
}
