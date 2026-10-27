package database

import (
	"gorm.io/gorm"

	"taskapp/internal/models"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Task{},
	)
}
