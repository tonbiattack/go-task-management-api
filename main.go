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
	taskHandler := handler.NewTaskHandler(taskRepo)

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

	// サーバーの起動
	router.Run(":8080")
}
