package storage

import (
	"errors"

	"taskapp/internal/models"
)

type MemoryStorage struct {
	users       map[uint]models.User
	userByEmail map[string]uint
	tasks       map[uint]models.Task

	nextUserID uint
	nextTaskID uint
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users:       make(map[uint]models.User),
		userByEmail: make(map[string]uint),
		tasks:       make(map[uint]models.Task),
		nextUserID:  1,
		nextTaskID:  1,
	}
}

var _ Storage = (*MemoryStorage)(nil)

func (s *MemoryStorage) CreateUser(user models.User) (uint, error) {
	if user.Email == "" {
		return 0, errors.New("email is empty")
	}

	if _, exists := s.userByEmail[user.Email]; exists {
		return 0, errors.New("email already exists")
	}

	user.ID = s.nextUserID
	s.nextUserID++

	s.users[user.ID] = user
	s.userByEmail[user.Email] = user.ID

	return user.ID, nil
}

func (s *MemoryStorage) GetUser(id uint) (models.User, error) {
	user, ok := s.users[id]
	if !ok {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}

func (s *MemoryStorage) GetUserByEmail(email string) (models.User, error) {
	id, ok := s.userByEmail[email]
	if !ok {
		return models.User{}, ErrUserNotFound
	}

	return s.GetUser(id)
}

func (s *MemoryStorage) CreateTask(
	userID uint,
	task models.Task,
) (uint, error) {
	if _, ok := s.users[userID]; !ok {
		return 0, ErrUserNotFound
	}
	if task.Title == "" {
		return 0, errors.New("title is empty")
	}

	task.ID = s.nextTaskID
	s.nextTaskID++
	task.UserID = userID

	s.tasks[task.ID] = task

	return task.ID, nil
}

func (s *MemoryStorage) GetTask(
	userID uint,
	taskID uint,
) (models.Task, error) {
	task, ok := s.tasks[taskID]
	if !ok || task.UserID != userID {
		return models.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (s *MemoryStorage) ListUserTasks(userID uint) ([]models.Task, error) {
	if _, ok := s.users[userID]; !ok {
		return nil, ErrUserNotFound
	}

	tasks := make([]models.Task, 0)

	for _, task := range s.tasks {
		if task.UserID == userID {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (s *MemoryStorage) UpdateTask(
	userID uint,
	taskID uint,
	update models.TaskUpdate,
) error {
	task, ok := s.tasks[taskID]
	if !ok || task.UserID != userID {
		return ErrTaskNotFound
	}

	if update.Title != nil {
		task.Title = *update.Title
	}
	if update.Tags != nil {
		task.Tags = *update.Tags
	}
	if update.Description != nil {
		task.Description = *update.Description
	}
	if update.IsCompleted != nil {
		task.IsCompleted = *update.IsCompleted
	}

	s.tasks[taskID] = task

	return nil
}

func (s *MemoryStorage) DeleteTask(userID uint, taskID uint) error {
	task, ok := s.tasks[taskID]
	if !ok || task.UserID != userID {
		return ErrTaskNotFound
	}

	delete(s.tasks, taskID)

	return nil
}
