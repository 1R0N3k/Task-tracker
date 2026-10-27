package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `gorm:"size:255;uniqueIndex;not null"`
	Name      string         `gorm:"size:100;not null"`
	Password  string         `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Tasks []Task `gorm:"foreignKey:UserID"`
}
