package main

import (
	"log"
	"os"

	"taskapp/internal/database"
	"taskapp/internal/models"
	"taskapp/internal/storage"
)

func main() {
	var store storage.Storage

	if os.Getenv("STORAGE") == "postgres" {
		db, err := database.Connect()
		if err != nil {
			log.Fatalf("connect: %v", err)
		}

		if err := database.Migrate(db); err != nil {
			log.Fatalf("migrate: %v", err)
		}

		store = storage.NewPostgresStorage(db)
	} else {
		store = storage.NewMemoryStorage()
	}

	userID, err := store.CreateUser(models.User{
		Email:    "alice@example.com",
		Name:     "Alice",
		Password: "hashed-password",
	})
	if err != nil {
		log.Fatalf("create user: %v", err)
	}

	taskID, err := store.CreateTask(userID, models.Task{
		Title:       "Изучить Go",
		Tags:        []string{"go", "learning"},
		Description: "Прочитать документацию",
	})
	if err != nil {
		log.Fatalf("create task: %v", err)
	}

	tasks, err := store.ListUserTasks(userID)
	if err != nil {
		log.Fatalf("list tasks: %v", err)
	}

	log.Printf("user %d has %d task(s), first task id: %d",
		userID, len(tasks), taskID)
}
