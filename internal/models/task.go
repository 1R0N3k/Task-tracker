package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;index"`
	Title       string `gorm:"size:255;not null"`
	Tags        []string
	Description string
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID"`
}

type TaskUpdate struct {
	Title       *string
	Tags        *[]string
	Description *string
	IsCompleted *bool
}
