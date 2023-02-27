package main

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	connectionString := "postgresql://postgres:password@localhost:5432/todos?sslmode=disable"

	e := echo.New()
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		e.Logger.Fatal(err)
	}

	repo := TodoRepository{db: db}
	h := Handler{todoRepo: repo}

	e.GET("/todos", h.list)
	e.POST("/todos", h.create)
	e.GET("/todos/:id", h.retrieve)
	e.POST("/todos/:id", h.update)
	e.DELETE("/todos/:id", h.delete)

	e.Logger.Fatal(e.Start(":1323"))
}
