package models

import (
	"time"

	"gorm.io/gorm"
)

type PomodoroSession struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	TaskID    *uint     // может быть null, если сессия не привязана к задаче
	Duration  int       `gorm:"not null"` // в минутах
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
}
