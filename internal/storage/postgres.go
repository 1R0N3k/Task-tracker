package storage

import (
	"errors"

	"gorm.io/gorm"

	"taskapp/internal/models"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(db *gorm.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

var _ Storage = (*PostgresStorage)(nil)

func (s *PostgresStorage) CreateUser(user models.User) (uint, error) {
	err := s.db.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *PostgresStorage) GetUser(id uint) (models.User, error) {
	var user models.User

	err := s.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *PostgresStorage) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := s.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *PostgresStorage) UpdateUser(id uint, user models.User) error {
	result := s.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"name":     user.Name,
			"email":    user.Email,
			"password": user.Password,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (s *PostgresStorage) DeleteUser(id uint) error {
	result := s.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (s *PostgresStorage) CreateTask(
	userID uint,
	task models.Task,
) (uint, error) {
	task.UserID = userID

	err := s.db.Create(&task).Error
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (s *PostgresStorage) GetTask(
	userID uint,
	taskID uint,
) (models.Task, error) {
	var task models.Task

	err := s.db.
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *PostgresStorage) ListUserTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task

	err := s.db.
		Where("user_id = ?", userID).
		Find(&tasks).Error

	return tasks, err
}

func (s *PostgresStorage) UpdateTask(
	userID uint,
	taskID uint,
	update models.TaskUpdate,
) error {
	updates := map[string]any{}

	if update.Title != nil {
		updates["title"] = *update.Title
	}
	if update.Tags != nil {
		updates["tags"] = *update.Tags
	}
	if update.Description != nil {
		updates["description"] = *update.Description
	}
	if update.IsCompleted != nil {
		updates["is_completed"] = *update.IsCompleted
	}

	result := s.db.Model(&models.Task{}).
		Where("id = ? AND user_id = ?", taskID, userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (s *PostgresStorage) DeleteTask(userID uint, taskID uint) error {
	result := s.db.
		Where("id = ? AND user_id = ?", taskID, userID).
		Delete(&models.Task{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
