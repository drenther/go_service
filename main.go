package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	scribble "github.com/nanobox-io/golang-scribble"
	uuid "github.com/satori/go.uuid"
)

type Task struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Body string `json:"body"`
}

var port = ":4021"
var dbName = "tasks"

var db scribble.Driver

func main() {
	d, dbErr := scribble.New(dbName, nil)
	if dbErr != nil {
		fmt.Println("Failed to connect to database", dbErr)
	}
	db = *d

	server := echo.New()

	server.GET("/", getTasks)

	server.Logger.Fatal(server.Start(port))
}

func createTask(context echo.Context) (err error) {
	user := strings.Split(strings.ToLower(context.Request().Header.Get("Authorization")), " ")[1]
	id := uuid.NewV4().String()
	task := &Task{
		Id:   id,
		User: user,
	}

	if err = db.Write(dbName, id, &task); err != nil {
		return
	}

	if err = context.Bind(task); err != nil {
		return
	}

	return context.JSON(http.StatusCreated, task)
}

func getTasks(context echo.Context) (err error) {
	records, err := db.ReadAll(dbName)
	if err != nil {
		return
	}

	tasks := []Task{}
	for _, f := range records {
		task := Task{}
		if err = json.Unmarshal([]byte(f), &task); err != nil {
			return
		}
		tasks = append(tasks, task)
	}

	return context.JSON(http.StatusOK, tasks)
}
