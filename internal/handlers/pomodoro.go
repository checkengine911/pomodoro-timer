package handlers

import (
	"net/http"
	"pomodoro-timer/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePomodoroSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var req struct {
			TaskID   *uint     `json:"task_id"`
			Duration int       `json:"duration" binding:"required"` // в минутах
			Start    time.Time `json:"start_time" binding:"required"`
			End      time.Time `json:"end_time" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		session := models.PomodoroSession{
			UserID:    userID,
			TaskID:    req.TaskID,
			Duration:  req.Duration,
			StartTime: req.Start,
			EndTime:   req.End,
		}
		if err := db.Create(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания сессии"})
			return
		}
		c.JSON(http.StatusOK, session)
	}
}

func GetPomodoroSessions(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var sessions []models.PomodoroSession
		if err := db.Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения сессий"})
			return
		}
		c.JSON(http.StatusOK, sessions)
	}
}

func UpdatePomodoroSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id := c.Param("id")
		var session models.PomodoroSession
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&session).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Сессия не найдена"})
			return
		}
		var req struct {
			TaskID   *uint      `json:"task_id"`
			Duration *int       `json:"duration"`
			Start    *time.Time `json:"start_time"`
			End      *time.Time `json:"end_time"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.TaskID != nil {
			session.TaskID = req.TaskID
		}
		if req.Duration != nil {
			session.Duration = *req.Duration
		}
		if req.Start != nil {
			session.StartTime = *req.Start
		}
		if req.End != nil {
			session.EndTime = *req.End
		}
		if err := db.Save(&session).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления сессии"})
			return
		}
		c.JSON(http.StatusOK, session)
	}
}

func DeletePomodoroSession(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		id := c.Param("id")
		if err := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.PomodoroSession{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления сессии"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}
}
