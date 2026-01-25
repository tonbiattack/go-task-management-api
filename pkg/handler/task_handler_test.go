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

func setupTaskHandler(t *testing.T) (*TaskHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	db := testutils.SetupTestDB(t)
	repo := repository.NewTaskRepository(db)
	handler := NewTaskHandler(repo)

	router := gin.New()
	return handler, router
}

func TestTaskHandler_CreateTask(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.POST("/task", handler.CreateTask)

	task := model.Task{
		Title:       "新しいタスク",
		Description: "新しいタスクの説明",
	}

	jsonData, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/task", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response model.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "新しいタスク", response.Title)
	assert.Equal(t, "新しいタスクの説明", response.Description)
	assert.NotEmpty(t, response.ID)
}

func TestTaskHandler_CreateTask_InvalidJSON(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.POST("/task", handler.CreateTask)

	req := httptest.NewRequest("POST", "/task", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.GET("/tasks", handler.GetAllTasks)

	// テストデータを作成
	testutils.CreateTestTask(handler.Repo.DB, "task-1", "タスク1", "説明1")
	testutils.CreateTestTask(handler.Repo.DB, "task-2", "タスク2", "説明2")

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*model.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestTaskHandler_GetTask(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.GET("/task/:id", handler.GetTask)

	// テストデータを作成
	testutils.CreateTestTask(handler.Repo.DB, "task-1", "テストタスク", "テスト説明")

	req := httptest.NewRequest("GET", "/task/task-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "task-1", response.ID)
	assert.Equal(t, "テストタスク", response.Title)
}

func TestTaskHandler_GetTask_NotFound(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.GET("/task/:id", handler.GetTask)

	req := httptest.NewRequest("GET", "/task/non-existent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Failed to retrieve task", response["error"])
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.PUT("/task/:id", handler.UpdateTask)

	// テストデータを作成
	testutils.CreateTestTask(handler.Repo.DB, "task-1", "元のタスク", "元の説明")

	updatedTask := model.Task{
		Title:       "更新されたタスク",
		Description: "更新された説明",
	}

	jsonData, _ := json.Marshal(updatedTask)
	req := httptest.NewRequest("PUT", "/task/task-1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "task-1", response.ID)
	assert.Equal(t, "更新されたタスク", response.Title)
	assert.Equal(t, "更新された説明", response.Description)
}

func TestTaskHandler_UpdateTask_InvalidJSON(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.PUT("/task/:id", handler.UpdateTask)

	req := httptest.NewRequest("PUT", "/task/task-1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	handler, router := setupTaskHandler(t)
	router.DELETE("/task/:id", handler.DeleteTask)

	// テストデータを作成
	testutils.CreateTestTask(handler.Repo.DB, "task-1", "削除されるタスク", "削除テスト")

	req := httptest.NewRequest("DELETE", "/task/task-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Task deleted successfully", response["message"])

	// タスクが削除されたことを確認
	_, err = handler.Repo.GetTaskByID("task-1")
	assert.Error(t, err)
}
