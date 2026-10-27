package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"taskapp/internal/models"
	"taskapp/internal/storage"
)

func buildApp(store storage.Storage) *echo.Echo {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello, Go!",
		})
	})

	e.POST("/users", func(c echo.Context) error {
		var user models.User

		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
		}

		id, err := store.CreateUser(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		user.ID = id

		return c.JSON(http.StatusCreated, user)
	})

	e.GET("/users/:id", func(c echo.Context) error {
		id, err := parseUint(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}

		user, err := store.GetUser(id)
		if errors.Is(err, storage.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, user)
	})

	e.POST("/users/:id/tasks", func(c echo.Context) error {
		userID, err := parseUint(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}

		var task models.Task

		if err := c.Bind(&task); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
		}

		taskID, err := store.CreateTask(userID, task)
		if errors.Is(err, storage.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		task.ID = taskID
		task.UserID = userID

		return c.JSON(http.StatusCreated, task)
	})

	e.GET("/users/:id/tasks", func(c echo.Context) error {
		userID, err := parseUint(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}

		tasks, err := store.ListUserTasks(userID)
		if errors.Is(err, storage.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, tasks)
	})

	e.GET("/users/:id/tasks/:taskID", func(c echo.Context) error {
		userID, err := parseUint(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}

		taskID, err := parseUint(c.Param("taskID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
		}

		task, err := store.GetTask(userID, taskID)
		if errors.Is(err, storage.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, task)
	})

	e.PATCH("/users/:id/tasks/:taskID", func(c echo.Context) error {
		userID, err := parseUint(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}

		taskID, err := parseUint(c.Param("taskID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
		}

		var update models.TaskUpdate

		if err := c.Bind(&update); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
		}

		if err := store.UpdateTask(userID, taskID, update); err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "updated",
		})
	})

	e.DELETE("/users/:id/tasks/:taskID", func(c echo.Context) error {
		userID, err := parseUint(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
		}

		taskID, err := parseUint(c.Param("taskID"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
		}

		if err := store.DeleteTask(userID, taskID); err != nil {
			if errors.Is(err, storage.ErrTaskNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "deleted",
		})
	})

	return e
}

func parseUint(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(parsed), nil
}
