package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tonbiattack/go-task-management-api/pkg/config"
	"github.com/tonbiattack/go-task-management-api/pkg/handler"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
)

func main() {
	db := config.InitDB()
	taskRepo := repository.NewTaskRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	workflowStatusRepo := repository.NewTaskWorkflowStatusRepository(db)
	workflowStepRepo := repository.NewWorkflowStepRepository(db)
	taskHandler := handler.NewTaskHandler(taskRepo)
	// ハンドラーの初期化
	workflowHandler := handler.NewWorkflowHandler(taskRepo, workflowRepo, workflowStatusRepo, workflowStepRepo)

	router := gin.Default()

	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// ルートの設定
	// /タスク関連のエンドポイント
	// GetAllTasks は、保存されているすべてのタスクを取得するエンドポイントです。
	router.GET("/tasks", taskHandler.GetAllTasks)

	// CreateTask は、新しいタスクを作成するエンドポイントです。
	router.POST("/task", taskHandler.CreateTask)

	// GetTask は、指定されたIDを持つ特定のタスクを取得するエンドポイントです。
	router.GET("/task/:id", taskHandler.GetTask)

	// UpdateTask は、指定されたIDのタスクを更新するエンドポイントです。
	router.PUT("/task/:id", taskHandler.UpdateTask)

	// DeleteTask は、指定されたIDのタスクを削除するエンドポイントです。
	router.DELETE("/task/:id", taskHandler.DeleteTask)

	//ワークフロー関連のエンドポイント
	// CreateWorkflow は、新しいワークフローを作成するエンドポイントです。
	router.POST("/workflows", workflowHandler.CreateWorkflow)

	// GetWorkflow は、指定されたIDを持つ特定のワークフローを取得するエンドポイントです。
	router.GET("/workflows/:id", workflowHandler.GetWorkflow)

	// AdvanceWorkflow は、指定されたタスクIDのワークフローステータスを次のステップに進めるエンドポイントです。
	// パスパラメーター: :task_id進行させたいタスクのid
	router.PUT("/task/advance/:task_id", workflowHandler.AdvanceWorkflow)

	// StartWorkflow は、特定のワークフローIDとタスクIDに関連付けてワークフローを開始するエンドポイントです。
	router.POST("/workflow/:workflow_id/start/:task_id", workflowHandler.StartWorkflow)

	// CreateWorkflowStep は、新しいワークフローステップを作成するエンドポイントです。
	router.POST("/workflow-steps", workflowHandler.CreateWorkflowStep)

	// サーバーの起動
	router.Run(":8080")
}
