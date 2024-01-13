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
// CreateWorkflow は新しいワークフローを作成するためのエンドポイントです。
// このメソッドは、クライアントからのJSON形式のワークフローデータを受け取り、
// データベースに保存します。
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
// GetWorkflow は特定のワークフローを取得するためのエンドポイントです。
// ワークフローIDはURLパラメータから取得されます。
func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	workflowID := c.Param("id")

	workflow, err := h.WorkflowRepo.FindByID(workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflow"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// AdvanceWorkflow は指定されたタスクIDに関連するワークフローステータスを
// 次のステップに進めるためのエンドポイントです。
// もし最終ステップに達している場合、ステータスは「COMPLETED」とマークされます。
func (h *WorkflowHandler) AdvanceWorkflow(c *gin.Context) {
	taskID := c.Param("task_id") // タスクIDの取得

	// 現在のワークフローステータスを取得
	currentStatus, err := h.WorkflowStatusRepo.FindByTaskID(taskID)
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

// StartWorkflow - POST /workflow/{workflow_id}/start/{task_id} に対するハンドラー関数
// StartWorkflow は指定されたワークフローIDとタスクIDに関連付けて
// ワークフローを開始するためのエンドポイントです。
// この操作により、初期ステップで新しいワークフローステータスが作成されます。
func (h *WorkflowHandler) StartWorkflow(c *gin.Context) {
	workflowID := c.Param("workflow_id") // ワークフローIDの取得
	taskID := c.Param("task_id")         // タスクIDの取得

	// ワークフローとタスクの存在を確認
	_, err := h.WorkflowRepo.FindByID(workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workflow not found"})
		return
	}
	_, err = h.TaskRepo.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Task not found"})
		return
	}

	// ワークフローの初期ステップを取得
	initialStep, err := h.WorkflowStepRepo.GetInitialStep(workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Initial workflow step not found"})
		return
	}

	// ワークフローステータスを初期ステップで作成
	status := model.TaskWorkflowStatus{
		ID:             uuid.New().String(),
		TaskID:         taskID,
		WorkflowStepID: initialStep.ID,
		Status:         "IN_PROGRESS",
	}
	if err := h.WorkflowStatusRepo.Create(&status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow started", "status": status})
}

// CreateWorkflowStep は新しいワークフローステップを作成するためのエンドポイントです。
// クライアントからのJSON形式のステップデータを受け取り、データベースに保存します。
func (h *WorkflowHandler) CreateWorkflowStep(c *gin.Context) {
	var step model.WorkflowStep
	if err := c.BindJSON(&step); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 新しいUUIDを生成してIDに設定
	step.ID = uuid.New().String()

	// ワークフローステップをデータベースに保存
	if err := h.WorkflowStepRepo.Create(&step); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow step"})
		return
	}

	c.JSON(http.StatusCreated, step)
}
