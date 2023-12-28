package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tonbiattack/go-task-management-api/pkg/config"
	"github.com/tonbiattack/go-task-management-api/pkg/handler"
	"github.com/tonbiattack/go-task-management-api/pkg/repository"
)

func main() {
	db := config.GetDB()
	taskRepo := repository.NewTaskRepository(db)
	taskHandler := handler.NewTaskHandler(taskRepo)
	router := mux.NewRouter()

	router.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	router.HandleFunc("/task", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/task/{id}", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/task/{id}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/task/{id}", taskHandler.DeleteTask).Methods("DELETE")

	corsMiddleware := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // viteアプリのホスト
	)

	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(router)); err != nil {
		log.Fatal(err)
	}
}
