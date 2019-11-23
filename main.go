package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	scribble "github.com/nanobox-io/golang-scribble"
	uuid "github.com/satori/go.uuid"
)

var port = ":4021"
var dbName = "tasks"

type ContextWithDB struct {
	echo.Context
	Db *scribble.Driver
}

func main() {
	db, dbErr := scribble.New(dbName, nil)
	if dbErr != nil {
		log.Fatal("Failed to connect to database", dbErr)
	}

	server := echo.New()

	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			context := &ContextWithDB{c, db}
			return next(context)
		}
	})

	server.GET("/", getTasks)

	server.POST("/", createTask)

	server.Logger.Fatal(server.Start(port))
}

type Task struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Body string `json:"body"`
}

func getUser(authHeader string) (string, error) {
	if authHeader == "" {
		return "", echo.ErrUnauthorized
	}

	substrings := strings.Split(strings.ToLower(authHeader), " ")

	if len(substrings) < 2 && substrings[0] != "" && substrings[1] != "" {
		return "", echo.ErrUnauthorized
	}

	return substrings[1], nil
}

func createTask(c echo.Context) (err error) {
	context := c.(*ContextWithDB)

	user, err := getUser(context.Request().Header.Get("Authorization"))
	if err != nil {
		return err
	}

	id := uuid.NewV4().String()
	task := &Task{
		Id:   id,
		User: user,
	}

	if err = context.Bind(task); err != nil {
		return
	}

	if err = context.Db.Write(dbName, id, &task); err != nil {
		return
	}

	return context.JSON(http.StatusCreated, task)
}

func getTasks(c echo.Context) (err error) {
	context := c.(*ContextWithDB)
	user, err := getUser(context.Request().Header.Get("Authorization"))
	if err != nil {
		return err
	}

	tasks := []Task{}
	records, err := context.Db.ReadAll(dbName)
	if err != nil {
		return context.JSON(http.StatusNoContent, tasks)
	}

	for _, f := range records {
		task := Task{}
		if err = json.Unmarshal([]byte(f), &task); err != nil {
			return
		}

		if user == task.User {
			tasks = append(tasks, task)
		}
	}

	return context.JSON(http.StatusOK, tasks)
}
