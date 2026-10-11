package main

import (
	"errors"
	"fmt"
)

type Task struct {
	Title       string
	Tags        []string
	Description string
	IsCompleted bool
}

type MemoryStorage struct {
	tasks map[string]Task
}

type TaskUpdate struct {
	Title       *string
	Tags        *[]string
	Description *string
	IsCompleted *bool
}

type Storage interface {
	CreateTask(
		title string,
		tags []string,
		description string,
		isCompleted bool,
	) error

	GetTask(title string) (Task, error)

	UpdateTask(
		title string,
		update TaskUpdate,
	) error

	DeleteTask(title string) error
}

func (storage *MemoryStorage) CreateTask(
	title string,
	tags []string,
	description string,
	isCompleted bool,
) error {
	if title == "" {
		return errors.New("title is empty")
	}

	storage.tasks[title] = Task{
		Title:       title,
		Tags:        tags,
		Description: description,
		IsCompleted: isCompleted,
	}

	return nil
}

func (storage *MemoryStorage) GetTask(title string) (Task, error) {
	task, ok := storage.tasks[title]
	if !ok {
		return Task{}, errors.New("task not found")
	}

	return task, nil
}

func (storage *MemoryStorage) UpdateTask(title string, update TaskUpdate) error {
	task, ok := storage.tasks[title]
	if !ok {
		return errors.New("task not found")
	}

	if update.Description != nil {
		task.Description = *update.Description
	}
	if update.Tags != nil {
		task.Tags = *update.Tags
	}
	if update.IsCompleted != nil {
		task.IsCompleted = *update.IsCompleted
	}

	if update.Title != nil {
		task.Title = *update.Title

		delete(storage.tasks, title)
		storage.tasks[*update.Title] = task
	} else {
		storage.tasks[title] = task
	}

	return nil
}

func (storage *MemoryStorage) DeleteTask(title string) error {
	delete(storage.tasks, title)
	return nil
}

func main() {
	var storage Storage = &MemoryStorage{
		tasks: make(map[string]Task),
	}

	// 1. CreateTask
	err := storage.CreateTask(
		"Задание",
		[]string{"homework", "Go lang"},
		"Описание задачи",
		false,
	)
	if err != nil {
		fmt.Println("Ошибка при создании задачи:", err)
		return
	}

	fmt.Println("Задача создана")

	// 2. GetTask
	task, err := storage.GetTask("Задание")
	if err != nil {
		fmt.Println("Ошибка при получении задачи:", err)
		return
	}

	fmt.Printf("Получена задача: %+v\n", task)

	// 3. UpdateTask
	description := "Новое описание"
	isCompleted := true

	err = storage.UpdateTask(
		"Задание",
		TaskUpdate{
			Description: &description,
			IsCompleted: &isCompleted,
		},
	)
	if err != nil {
		fmt.Println("Ошибка при обновлении задачи:", err)
		return
	}

	fmt.Println("Задача обновлена")

	// Проверяем, что обновление сработало
	updatedTask, err := storage.GetTask("Задание")
	if err != nil {
		fmt.Println("Ошибка при получении обновлённой задачи:", err)
		return
	}

	fmt.Printf("Задача после обновления: %+v\n", updatedTask)

	// 4. DeleteTask
	err = storage.DeleteTask("Задание")
	if err != nil {
		fmt.Println("Ошибка при удалении задачи:", err)
		return
	}

	fmt.Println("Задача удалена")

	// Проверяем, что задачи больше нет
	_, err = storage.GetTask("Задание")
	if err != nil {
		fmt.Println("Ожидаемая ошибка после удаления:", err)
	} else {
		fmt.Println("Ошибка: задача всё ещё существует")
	}
}
