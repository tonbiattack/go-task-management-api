package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tonbiattack/go-task-management-api/pkg/config"
	"github.com/tonbiattack/go-task-management-api/pkg/handler"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
)

func main() {
	// データベース接続を取得
	db := config.GetDB()

	// リポジトリとハンドラーをセットアップ
	taskRepo := repository.NewTaskRepository(db)
	taskHandler := handler.NewTaskHandler(taskRepo)

	// Gorilla Muxルーターを作成
	router := mux.NewRouter()

	// エンドポイントとハンドラーを紐付け
	router.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	router.HandleFunc("/task", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/task/{id}", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/task/{id}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/task/{id}", taskHandler.DeleteTask).Methods("DELETE") // 特定のタスクを削除するエンドポイント

	// HTTPサーバーを8080ポートで起動
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
