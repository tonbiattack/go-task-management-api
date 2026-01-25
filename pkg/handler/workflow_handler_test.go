package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
	"github.com/tonbiattack/go-task-management-api/pkg/testutils"
)

func setupWorkflowHandler(t *testing.T) (*WorkflowHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	db := testutils.SetupTestDB(t)
	taskRepo := repository.NewTaskRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	workflowStatusRepo := repository.NewTaskWorkflowStatusRepository(db)
	workflowStepRepo := repository.NewWorkflowStepRepository(db)

	handler := NewWorkflowHandler(taskRepo, workflowRepo, workflowStatusRepo, workflowStepRepo)

	router := gin.New()
	return handler, router
}

func TestWorkflowHandler_CreateWorkflow(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.POST("/workflows", handler.CreateWorkflow)

	workflow := model.Workflow{
		Name:        "新しいワークフロー",
		Description: "新しいワークフローの説明",
	}

	jsonData, _ := json.Marshal(workflow)
	req := httptest.NewRequest("POST", "/workflows", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Workflow
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "新しいワークフロー", response.Name)
	assert.Equal(t, "新しいワークフローの説明", response.Description)
	assert.NotEmpty(t, response.ID)
}

func TestWorkflowHandler_CreateWorkflow_InvalidJSON(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.POST("/workflows", handler.CreateWorkflow)

	req := httptest.NewRequest("POST", "/workflows", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestWorkflowHandler_GetWorkflow(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.GET("/workflows/:id", handler.GetWorkflow)

	// テストデータを作成
	db := testutils.SetupTestDB(t)
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "テスト説明")

	req := httptest.NewRequest("GET", "/workflows/workflow-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Workflow
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "workflow-1", response.ID)
	assert.Equal(t, "テストワークフロー", response.Name)
}

func TestWorkflowHandler_GetWorkflow_NotFound(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.GET("/workflows/:id", handler.GetWorkflow)

	req := httptest.NewRequest("GET", "/workflows/non-existent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Failed to retrieve workflow", response["error"])
}

func TestWorkflowHandler_StartWorkflow(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.POST("/workflow/:workflow_id/start/:task_id", handler.StartWorkflow)

	// テストデータを作成
	db := handler.TaskRepo.DB
	testutils.CreateTestTask(db, "task-1", "テストタスク", "テスト説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "テスト説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)

	req := httptest.NewRequest("POST", "/workflow/workflow-1/start/task-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Workflow started", response["message"])
	assert.NotNil(t, response["status"])
}

func TestWorkflowHandler_StartWorkflow_WorkflowNotFound(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.POST("/workflow/:workflow_id/start/:task_id", handler.StartWorkflow)

	// タスクのみ作成（ワークフローは作成しない）
	db := handler.TaskRepo.DB
	testutils.CreateTestTask(db, "task-1", "テストタスク", "テスト説明")

	req := httptest.NewRequest("POST", "/workflow/non-existent/start/task-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Workflow not found", response["error"])
}

func TestWorkflowHandler_StartWorkflow_TaskNotFound(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.POST("/workflow/:workflow_id/start/:task_id", handler.StartWorkflow)

	// ワークフローのみ作成（タスクは作成しない）
	db := handler.TaskRepo.DB
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "テスト説明")

	req := httptest.NewRequest("POST", "/workflow/workflow-1/start/non-existent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Task not found", response["error"])
}

func TestWorkflowHandler_AdvanceWorkflow(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.PUT("/task/advance/:task_id", handler.AdvanceWorkflow)

	// テストデータを作成
	db := handler.TaskRepo.DB
	testutils.CreateTestTask(db, "task-1", "テストタスク", "テスト説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "テスト説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestWorkflowStep(db, "step-2", "workflow-1", "ステップ2", 2)
	testutils.CreateTestTaskWorkflowStatus(db, "status-1", "task-1", "step-1", "IN_PROGRESS")

	req := httptest.NewRequest("PUT", "/task/advance/task-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Workflow advanced to the next step", response["message"])
}

func TestWorkflowHandler_AdvanceWorkflow_FinalStep(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.PUT("/task/advance/:task_id", handler.AdvanceWorkflow)

	// テストデータを作成（最終ステップに設定）
	db := handler.TaskRepo.DB
	testutils.CreateTestTask(db, "task-1", "テストタスク", "テスト説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "テスト説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestTaskWorkflowStatus(db, "status-1", "task-1", "step-1", "IN_PROGRESS")

	req := httptest.NewRequest("PUT", "/task/advance/task-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Workflow completed", response["message"])
}

func TestWorkflowHandler_CreateWorkflowStep(t *testing.T) {
	handler, router := setupWorkflowHandler(t)
	router.POST("/workflow-steps", handler.CreateWorkflowStep)

	// 前提条件としてワークフローを作成
	db := handler.TaskRepo.DB
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "テスト説明")

	step := model.WorkflowStep{
		WorkflowID:  "workflow-1",
		Name:        "新しいステップ",
		Description: "新しいステップの説明",
		Order:       1,
	}

	jsonData, _ := json.Marshal(step)
	req := httptest.NewRequest("POST", "/workflow-steps", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.WorkflowStep
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "workflow-1", response.WorkflowID)
	assert.Equal(t, "新しいステップ", response.Name)
	assert.Equal(t, 1, response.Order)
	assert.NotEmpty(t, response.ID)
}
