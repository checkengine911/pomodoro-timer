package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskAnalytics struct {
	TaskID   uint   `json:"task_id"`
	Title    string `json:"title"`
	Duration int    `json:"total_minutes"`
}

func TasksTimeAnalytics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var results []TaskAnalytics
		rows, err := db.Raw(`
			SELECT t.id as task_id, t.title, COALESCE(SUM(p.duration),0) as duration
			FROM tasks t
			LEFT JOIN pomodoro_sessions p ON t.id = p.task_id AND p.user_id = ?
			WHERE t.user_id = ?
			GROUP BY t.id, t.title
		`, userID, userID).Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка аналитики"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var a TaskAnalytics
			if err := rows.Scan(&a.TaskID, &a.Title, &a.Duration); err == nil {
				results = append(results, a)
			}
		}
		c.JSON(http.StatusOK, results)
	}
}
