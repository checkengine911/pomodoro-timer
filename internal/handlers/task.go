package handlers

import (
	"net/http"
	"pomodoro-timer/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var req struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task := models.Task{
			UserID:      userID,
			Title:       req.Title,
			Description: req.Description,
			Status:      "pending",
		}
		if err := db.Create(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания задачи"})
			return
		}
		c.JSON(http.StatusOK, task)
	}
}

func GetTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var tasks []models.Task
		if err := db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения задач"})
			return
		}
		c.JSON(http.StatusOK, tasks)
	}
}

func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id := c.Param("id")
		var task models.Task
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}
		var req struct {
			Title       *string `json:"title"`
			Description *string `json:"description"`
			Status      *string `json:"status"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Title != nil {
			task.Title = *req.Title
		}
		if req.Description != nil {
			task.Description = *req.Description
		}
		if req.Status != nil {
			task.Status = *req.Status
		}
		if err := db.Save(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления задачи"})
			return
		}
		c.JSON(http.StatusOK, task)
	}
}

func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id := c.Param("id")
		if err := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления задачи"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}
}
