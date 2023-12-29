package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
)

type TaskHandler struct {
	Repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

// CreateTask - POST /tasks に対するハンドラー関数です。
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task model.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task.ID = uuid.New().String()

	if err := h.Repo.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	savedTask, err := h.Repo.GetTaskByID(task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve saved task"})
		return
	}

	c.JSON(http.StatusCreated, savedTask)
}

// GET /tasks に対するハンドラー関数です。
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.Repo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask - GET /tasks/{id} に対するハンドラー関数です。
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")

	task, err := h.Repo.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask - PUT /tasks/{id} に対するハンドラー関数です。
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")

	var task model.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	task.ID = taskID

	if err := h.Repo.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	updatedTask, err := h.Repo.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask - DELETE /tasks/{id} に対するハンドラー関数です。
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	if err := h.Repo.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.Status(http.StatusNoContent)
}
