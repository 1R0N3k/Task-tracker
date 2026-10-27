package storage

import (
	"errors"

	"taskapp/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTaskNotFound = errors.New("task not found")
)

type Storage interface {
	CreateUser(user models.User) (uint, error)
	GetUser(id uint) (models.User, error)
	GetUserByEmail(email string) (models.User, error)

	CreateTask(userID uint, task models.Task) (uint, error)
	GetTask(userID, taskID uint) (models.Task, error)
	ListUserTasks(userID uint) ([]models.Task, error)
	UpdateTask(userID, taskID uint, update models.TaskUpdate) error
	DeleteTask(userID, taskID uint) error
}
