package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
)

type WorkflowHandler struct {
	TaskRepo           *repository.TaskRepository
	WorkflowRepo       *repository.WorkflowRepository
	WorkflowStatusRepo *repository.TaskWorkflowStatusRepository
	WorkflowStepRepo   *repository.WorkflowStepRepository
}

func NewWorkflowHandler(
	taskRepo *repository.TaskRepository,
	workflowRepo *repository.WorkflowRepository,
	workflowStatusRepo *repository.TaskWorkflowStatusRepository,
	workflowStepRepo *repository.WorkflowStepRepository,
) *WorkflowHandler {
	return &WorkflowHandler{
		TaskRepo:           taskRepo,
		WorkflowRepo:       workflowRepo,
		WorkflowStatusRepo: workflowStatusRepo,
		WorkflowStepRepo:   workflowStepRepo,
	}
}

// CreateWorkflow - POST /workflows に対するハンドラー関数
func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	var workflow model.Workflow
	if err := c.BindJSON(&workflow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	workflow.ID = uuid.New().String()

	if err := h.WorkflowRepo.Create(&workflow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow"})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

// GetWorkflow - GET /workflows/{id} に対するハンドラー関数
func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	workflowID := c.Param("id")

	workflow, err := h.WorkflowRepo.FindByID(workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) AdvanceWorkflow(c *gin.Context) {
	taskID := c.Param("task_id") // タスクIDの取得

	// 現在のワークフローステータスを取得
	currentStatus, err := h.WorkflowStatusRepo.FindByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve current workflow status"})
		return
	}

	// 次のワークフローステップを取得
	nextStep, err := h.WorkflowStepRepo.GetNextStep(currentStatus.WorkflowStepID)
	if err != nil {
		if errors.Is(err, repository.ErrNoNextStep) {
			// 最終ステップに到達した場合の処理
			currentStatus.Status = "COMPLETED"
			if err := h.WorkflowStatusRepo.Update(currentStatus); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark workflow as completed"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Workflow completed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve next workflow step"})
		return
	}

	// 次のステップに進む通常の処理
	currentStatus.WorkflowStepID = nextStep.ID
	if err := h.WorkflowStatusRepo.Update(currentStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow advanced to the next step"})
}
