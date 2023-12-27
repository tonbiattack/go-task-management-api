package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
)

// TaskHandler - タスク関連のHTTPハンドラーを保持する構造体です。
type TaskHandler struct {
	Repo *repository.TaskRepository
}

// NewTaskHandler - 新しいTaskHandlerを作成します。
func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

// CreateTask - POST /tasks に対するハンドラー関数です。
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// UUID生成
	task.ID = uuid.New().String()

	// データベースにタスクを保存
	if err := h.Repo.CreateTask(&task); err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// 保存されたタスクの詳細を取得
	savedTask, err := h.Repo.GetTaskByID(task.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve saved task", http.StatusInternalServerError)
		return
	}

	// 保存されたタスクの詳細をJSONとしてレスポンスに書き込む
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedTask)
}

// GET /tasks に対するハンドラー関数です。
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Repo.GetAllTasks()
	if err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	// 取得したタスクをJSONとしてレスポンスに書き込む
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// GetTask - GET /tasks/{id} に対するハンドラー関数です。
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// URLからタスクIDを取得
	vars := mux.Vars(r)
	taskID := vars["id"]

	// データベースからタスクを取得
	task, err := h.Repo.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		return
	}

	// 取得したタスクの詳細をJSONとしてレスポンスに書き込む
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// UpdateTask - PUT /tasks/{id} に対するハンドラー関数です。
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// URLからタスクIDを取得
	vars := mux.Vars(r)
	taskID := vars["id"]

	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	task.ID = taskID

	// データベースにタスクを更新
	if err := h.Repo.UpdateTask(&task); err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// 更新されたタスクの詳細を取得
	updatedTask, err := h.Repo.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Failed to retrieve updated task", http.StatusInternalServerError)
		return
	}

	// 更新されたタスクの詳細をJSONとしてレスポンスに書き込む
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

// DeleteTask - DELETE /tasks/{id} に対するハンドラー関数です。
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// URLからタスクIDを取得
	vars := mux.Vars(r)
	taskID := vars["id"]

	// データベースからタスクを削除
	if err := h.Repo.DeleteTask(taskID); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// 削除成功のレスポンスを返す
	w.WriteHeader(http.StatusNoContent)
}
