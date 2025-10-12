package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&model.Task{}, &model.Workflow{}, &model.WorkflowStep{}, &model.TaskWorkflowStatus{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func seedWorkflowData(t *testing.T, db *gorm.DB) *model.Task {
	t.Helper()

	now := time.Now().UTC()
	status := "IN_PROGRESS"

	task := &model.Task{
		ID:          "task-join-1",
		Title:       "Design database schema",
		Description: "Create tables and relations",
		Status:      &status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	workflow := &model.Workflow{
		ID:          "workflow-1",
		Name:        "Development Flow",
		Description: "Standard software delivery workflow",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	step := &model.WorkflowStep{
		ID:          "workflow-step-1",
		WorkflowID:  workflow.ID,
		Name:        "Design",
		Description: "Design deliverables preparation",
		Order:       1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskWorkflowStatus := &model.TaskWorkflowStatus{
		ID:             "task-workflow-status-1",
		TaskID:         task.ID,
		WorkflowStepID: step.ID,
		Status:         status,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := db.Create(workflow).Error; err != nil {
		t.Fatalf("failed to create workflow: %v", err)
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	if err := db.Create(step).Error; err != nil {
		t.Fatalf("failed to create workflow step: %v", err)
	}
	if err := db.Create(taskWorkflowStatus).Error; err != nil {
		t.Fatalf("failed to create task workflow status: %v", err)
	}

	return task
}

func TestTaskRepository_GetTaskWorkflowDetail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)
	task := seedWorkflowData(t, db)

	detail, err := repo.GetTaskWorkflowDetail(task.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if detail.TaskID != task.ID {
		t.Errorf("expected task id %s, got %s", task.ID, detail.TaskID)
	}
	if detail.WorkflowName != "Development Flow" {
		t.Errorf("expected workflow name Development Flow, got %s", detail.WorkflowName)
	}
	if detail.WorkflowStepName != "Design" {
		t.Errorf("expected workflow step name Design, got %s", detail.WorkflowStepName)
	}
	if detail.WorkflowStatus != "IN_PROGRESS" {
		t.Errorf("expected workflow status IN_PROGRESS, got %s", detail.WorkflowStatus)
	}
	if detail.TaskStatus == nil || *detail.TaskStatus != "IN_PROGRESS" {
		t.Errorf("expected task status pointer to IN_PROGRESS, got %v", detail.TaskStatus)
	}
}

func TestTaskRepository_GetTaskWorkflowDetail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	_, err := repo.GetTaskWorkflowDetail("unknown-task")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got %v", err)
	}
}
