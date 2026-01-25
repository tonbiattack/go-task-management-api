package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonbiattack/go-task-management-api/pkg/config"
	"github.com/tonbiattack/go-task-management-api/pkg/handler"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
)

func setupTestApp(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	db := config.InitTestDB()

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})

	taskRepo := repository.NewTaskRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	workflowStatusRepo := repository.NewTaskWorkflowStatusRepository(db)
	workflowStepRepo := repository.NewWorkflowStepRepository(db)

	taskHandler := handler.NewTaskHandler(taskRepo)
	workflowHandler := handler.NewWorkflowHandler(taskRepo, workflowRepo, workflowStatusRepo, workflowStepRepo)

	router := gin.New()

	// タスク関連のエンドポイント
	router.GET("/tasks", taskHandler.GetAllTasks)
	router.POST("/task", taskHandler.CreateTask)
	router.GET("/task/:id", taskHandler.GetTask)
	router.PUT("/task/:id", taskHandler.UpdateTask)
	router.DELETE("/task/:id", taskHandler.DeleteTask)

	// ワークフロー関連のエンドポイント
	router.POST("/workflows", workflowHandler.CreateWorkflow)
	router.GET("/workflows/:id", workflowHandler.GetWorkflow)
	router.PUT("/task/advance/:task_id", workflowHandler.AdvanceWorkflow)
	router.POST("/workflow/:workflow_id/start/:task_id", workflowHandler.StartWorkflow)
	router.POST("/workflow-steps", workflowHandler.CreateWorkflowStep)

	return router
}

func TestTaskWorkflowIntegration(t *testing.T) {
	app := setupTestApp(t)

	// 1. タスクを作成
	task := model.Task{
		Title:       "統合テストタスク",
		Description: "統合テスト用のタスクです",
	}

	jsonData, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/task", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdTask model.Task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	require.NoError(t, err)
	taskID := createdTask.ID

	// 2. ワークフローを作成
	workflow := model.Workflow{
		Name:        "統合テストワークフロー",
		Description: "統合テスト用のワークフローです",
	}

	jsonData, _ = json.Marshal(workflow)
	req = httptest.NewRequest("POST", "/workflows", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdWorkflow model.Workflow
	err = json.Unmarshal(w.Body.Bytes(), &createdWorkflow)
	require.NoError(t, err)
	workflowID := createdWorkflow.ID

	// 3. ワークフローステップを作成
	step1 := model.WorkflowStep{
		WorkflowID:  workflowID,
		Name:        "企画",
		Description: "企画段階",
		Order:       1,
	}

	jsonData, _ = json.Marshal(step1)
	req = httptest.NewRequest("POST", "/workflow-steps", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	step2 := model.WorkflowStep{
		WorkflowID:  workflowID,
		Name:        "実行",
		Description: "実行段階",
		Order:       2,
	}

	jsonData, _ = json.Marshal(step2)
	req = httptest.NewRequest("POST", "/workflow-steps", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	step3 := model.WorkflowStep{
		WorkflowID:  workflowID,
		Name:        "完了",
		Description: "完了段階",
		Order:       3,
	}

	jsonData, _ = json.Marshal(step3)
	req = httptest.NewRequest("POST", "/workflow-steps", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	// 4. ワークフローを開始
	req = httptest.NewRequest("POST", "/workflow/"+workflowID+"/start/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var startResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &startResponse)
	require.NoError(t, err)
	assert.Equal(t, "Workflow started", startResponse["message"])

	// 5. ワークフローを次のステップに進める（企画 → 実行）
	req = httptest.NewRequest("PUT", "/task/advance/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var advanceResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &advanceResponse)
	require.NoError(t, err)
	assert.Equal(t, "Workflow advanced to the next step", advanceResponse["message"])

	// 6. ワークフローを次のステップに進める（実行 → 完了）
	req = httptest.NewRequest("PUT", "/task/advance/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &advanceResponse)
	require.NoError(t, err)
	assert.Equal(t, "Workflow advanced to the next step", advanceResponse["message"])

	// 7. 最終ステップで完了させる
	req = httptest.NewRequest("PUT", "/task/advance/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &advanceResponse)
	require.NoError(t, err)
	assert.Equal(t, "Workflow completed", advanceResponse["message"])

	// 8. タスクが正常に取得できることを確認
	req = httptest.NewRequest("GET", "/task/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var retrievedTask model.Task
	err = json.Unmarshal(w.Body.Bytes(), &retrievedTask)
	require.NoError(t, err)
	assert.Equal(t, taskID, retrievedTask.ID)
	assert.Equal(t, "統合テストタスク", retrievedTask.Title)
}

func TestTaskCRUDOperations(t *testing.T) {
	app := setupTestApp(t)

	// 1. タスクを作成
	task := model.Task{
		Title:       "CRUDテストタスク",
		Description: "CRUD操作のテストです",
	}

	jsonData, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/task", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var createdTask model.Task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	require.NoError(t, err)
	taskID := createdTask.ID

	// 2. 作成したタスクを取得
	req = httptest.NewRequest("GET", "/task/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var retrievedTask model.Task
	err = json.Unmarshal(w.Body.Bytes(), &retrievedTask)
	require.NoError(t, err)
	assert.Equal(t, taskID, retrievedTask.ID)
	assert.Equal(t, "CRUDテストタスク", retrievedTask.Title)

	// 3. タスクを更新
	updatedTask := model.Task{
		Title:       "更新されたCRUDテストタスク",
		Description: "更新されたCRUD操作のテストです",
	}

	jsonData, _ = json.Marshal(updatedTask)
	req = httptest.NewRequest("PUT", "/task/"+taskID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var updatedResult model.Task
	err = json.Unmarshal(w.Body.Bytes(), &updatedResult)
	require.NoError(t, err)
	assert.Equal(t, "更新されたCRUDテストタスク", updatedResult.Title)

	// 4. すべてのタスクを取得
	req = httptest.NewRequest("GET", "/tasks", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var allTasks []*model.Task
	err = json.Unmarshal(w.Body.Bytes(), &allTasks)
	require.NoError(t, err)
	assert.Len(t, allTasks, 1)
	assert.Equal(t, taskID, allTasks[0].ID)

	// 5. タスクを削除
	req = httptest.NewRequest("DELETE", "/task/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// 6. 削除されたタスクが取得できないことを確認
	req = httptest.NewRequest("GET", "/task/"+taskID, nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
