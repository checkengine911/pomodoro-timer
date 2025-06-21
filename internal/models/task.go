package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	Title       string `gorm:"not null"`
	Description string
	Status      string `gorm:"default:'pending'"` // pending, in_progress, done
}
