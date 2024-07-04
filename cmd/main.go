package main

import (
	"log"

	"github.com/Egorpalan/time-tracker/internal/handlers"
	"github.com/Egorpalan/time-tracker/internal/repository"
	"github.com/Egorpalan/time-tracker/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	repository.InitDB(cfg)

	r := gin.Default()

	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.AddUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)
	r.GET("/users/:id/tasks", handlers.GetUserTasks)
	r.POST("/tasks", handlers.StartTask)
	r.PUT("/tasks/:task_id/end", handlers.EndTask)

	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
