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

	e.GET("/todos", func(c echo.Context) error {
		return indexHandler(c, db)
	})

	e.POST("/todos", func(c echo.Context) error {
		return createHandler(c, db)
	})

	e.GET("/todos/:id", func(c echo.Context) error {
		return retrieveHandler(c, db)
	})

	e.POST("/todos/:id", func(c echo.Context) error {
		return updateHandler(c, db)
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		return deleteHandler(c, db)
	})

	e.Logger.Fatal(e.Start(":1323"))

}
