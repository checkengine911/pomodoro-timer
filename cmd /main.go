package main

import (
	"log"
	"os"
	"pomodoro-timer/internal/handlers"
	"pomodoro-timer/internal/middleware"
	"pomodoro-timer/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=pomodoro port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Task{}, &models.PomodoroSession{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	r := gin.Default()

	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))

	taskGroup := r.Group("/tasks")
	taskGroup.Use(middleware.JWTAuth())
	taskGroup.POST("", handlers.CreateTask(db))
	taskGroup.GET("", handlers.GetTasks(db))
	taskGroup.PUT(":id", handlers.UpdateTask(db))
	taskGroup.DELETE(":id", handlers.DeleteTask(db))

	pomodoroGroup := r.Group("/pomodoro")
	pomodoroGroup.Use(middleware.JWTAuth())
	pomodoroGroup.POST("", handlers.CreatePomodoroSession(db))
	pomodoroGroup.GET("", handlers.GetPomodoroSessions(db))
	pomodoroGroup.PUT(":id", handlers.UpdatePomodoroSession(db))
	pomodoroGroup.DELETE(":id", handlers.DeletePomodoroSession(db))

	r.GET("/analytics/tasks-time", middleware.JWTAuth(), handlers.TasksTimeAnalytics(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
