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
	router.GET("/tasks", taskHandler.GetAllTasks)
	router.POST("/task", taskHandler.CreateTask)
	router.GET("/task/:id", taskHandler.GetTask)
	router.PUT("/task/:id", taskHandler.UpdateTask)
	router.DELETE("/task/:id", taskHandler.DeleteTask)

	router.POST("/workflows", workflowHandler.CreateWorkflow)
	router.GET("/workflows/:id", workflowHandler.GetWorkflow)
	router.PUT("/task/advance/:task_id", workflowHandler.AdvanceWorkflow)

	// サーバーの起動
	router.Run(":8080")
}
