package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Handler struct {
	db *sql.DB
}

func (h Handler) getTodo(id int) (Todo, error) {
	var res Todo

	rows, err := h.db.Query("SELECT * FROM todos WHERE id=$1", id)
	if err != nil {
		return res, err
	}
	if !rows.Next() {
		return res, echo.ErrNotFound
	}
	if err := rows.Scan(&res.ID, &res.Description, &res.IsDone); err != nil {
		return res, err
	}

	return res, nil
}

func (h Handler) indexHandler(c echo.Context) error {
	var todos []Todo

	log := c.Logger()

	rows, err := h.db.Query("SELECT * FROM todos")
	if err != nil {
		return c.String(http.StatusFailedDependency, "An error occured.")
	}

	for rows.Next() {
		var res Todo
		if err := rows.Scan(&res.ID, &res.Description, &res.IsDone); err != nil {
			log.Error("An error occured while executing query: %v", err)
		}
		todos = append(todos, res)
	}

	return c.JSON(http.StatusOK, todos)
}

func (h Handler) createHandler(c echo.Context) error {
	var res Todo
	var id int
	description := c.FormValue("description")
	log := c.Logger()

	// result, err := db.Exec("INSERT INTO todos (description,is_done) VALUES ('hvhvhvh','f')  RETURNING id ", description)
	err := h.db.QueryRow("INSERT INTO todos ( description, is_done ) VALUES ( $1, false ) RETURNING id", description).Scan(&id)
	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}

	res, err = h.getTodo(id)

	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) deleteHandler(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "ID must be integer.")
	}

	log := c.Logger()

	_, err = h.db.Exec("DELETE from todos WHERE id=$1", id)
	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}
	return c.String(http.StatusOK, "")
}

func (h Handler) retrieveHandler(c echo.Context) error {
	log := c.Logger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "ID must be integer.")
	}

	log.Error("get todo", id)

	todo, err := h.getTodo(id)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, todo)
}

func (h Handler) updateHandler(c echo.Context) error {
	log := c.Logger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "ID must be integer.")
	}

	description := c.FormValue("description")
	isDone := c.FormValue("is_done")
	result, err := h.db.Exec("UPDATE todos SET description=$1, is_done=$2 WHERE id=$3", description, isDone, id)
	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		log.Fatal("An error occured while executing query: %v", err)
	}

	todo, err := h.getTodo(id)
	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}
	return c.JSON(http.StatusOK, todo)
}
