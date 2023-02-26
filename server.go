package main

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID          int64  `json:"id" xml:"id"`
	Description string `json:"description" xml:"description"`
	IsDone      bool   `json:"is_done" xml:"isDone"`
}

func main() {
	connectionString := "postgresql://postgres:password@localhost:5432/todos?sslmode=disable"

	e := echo.New()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		e.Logger.Fatal(err)
	}
	h := Handler{db: db}

	e.GET("/todos", h.indexHandler)
	e.POST("/todos", h.createHandler)
	e.GET("/todos/:id", h.retrieveHandler)
	e.POST("/todos/:id", h.updateHandler)
	e.DELETE("/todos/:id", h.deleteHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
