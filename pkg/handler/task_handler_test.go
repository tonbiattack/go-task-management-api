package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHandlerTestDB(t *testing.T) *gorm.DB {
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

func seedHandlerWorkflowData(t *testing.T, db *gorm.DB) *model.Task {
	t.Helper()

	now := time.Now().UTC()
	status := "IN_PROGRESS"

	task := &model.Task{
		ID:          "task-handler-1",
		Title:       "Implement workflow join",
		Description: "Ensure task joins workflow data",
		Status:      &status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	workflow := &model.Workflow{
		ID:          "workflow-handler-1",
		Name:        "Sample Workflow",
		Description: "Workflow used for join testing",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	step := &model.WorkflowStep{
		ID:          "workflow-step-handler-1",
		WorkflowID:  workflow.ID,
		Name:        "Draft",
		Description: "Drafting specifications",
		Order:       1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskWorkflowStatus := &model.TaskWorkflowStatus{
		ID:             "task-workflow-status-handler-1",
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

func TestTaskHandler_GetTaskWorkflowDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupHandlerTestDB(t)
	task := seedHandlerWorkflowData(t, db)

	taskHandler := NewTaskHandler(repository.NewTaskRepository(db))
	router := gin.New()
	router.GET("/task/:id/workflow-detail", taskHandler.GetTaskWorkflowDetail)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s/workflow-detail", task.ID), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var detail model.TaskWorkflowDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if detail.TaskID != task.ID {
		t.Errorf("expected task id %s, got %s", task.ID, detail.TaskID)
	}
	if detail.WorkflowStepName != "Draft" {
		t.Errorf("expected step name Draft, got %s", detail.WorkflowStepName)
	}
	if detail.WorkflowName != "Sample Workflow" {
		t.Errorf("expected workflow name Sample Workflow, got %s", detail.WorkflowName)
	}
}

func TestTaskHandler_GetTaskWorkflowDetail_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupHandlerTestDB(t)
	taskHandler := NewTaskHandler(repository.NewTaskRepository(db))
	router := gin.New()
	router.GET("/task/:id/workflow-detail", taskHandler.GetTaskWorkflowDetail)

	req := httptest.NewRequest(http.MethodGet, "/task/unknown/workflow-detail", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
