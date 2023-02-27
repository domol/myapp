package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Handler struct {
	todoRepo TodoRepository
}

func (h Handler) list(c echo.Context) error {
	log := c.Logger()

	todos, err := h.todoRepo.list()
	if err != nil {
		log.Error("An error occured while executing query: %v", err)
	}

	return c.JSON(http.StatusOK, todos)
}

func (h Handler) create(c echo.Context) error {
	var res Todo
	description := c.FormValue("description")
	if len(description) == 0 {
		return c.String(http.StatusExpectationFailed, "Field description is required.")
	}
	log := c.Logger()
	res, err := h.todoRepo.create(description)
	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "ID must be integer.")
	}

	log := c.Logger()

	err = h.todoRepo.delete(int64(id))
	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}

	return c.String(http.StatusOK, "")
}

func (h Handler) retrieve(c echo.Context) error {
	log := c.Logger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "ID must be integer.")
	}

	log.Error("get todo", id)

	todo, err := h.todoRepo.get(int64(id))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, todo)
}

func (h Handler) _getFormValues(c echo.Context, todo *Todo) (err error) {
	todo.Description = c.FormValue("description")

	isDoneString := c.FormValue("is_done")
	todo.IsDone, err = strconv.ParseBool(isDoneString)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "is_done must be true or false.")
	}
	return
}

func (h Handler) update(c echo.Context) error {
	log := c.Logger()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "ID must be integer.")
	}

	data := Todo{}
	h._getFormValues(c, &data)

	err = h.todoRepo.update(int64(id), data)
	if err != nil {
		log.Fatal("Could not update todo. %v", err)
	}

	todo, err := h.todoRepo.get(int64(id))
	if err != nil {
		log.Fatal("An error occured while executing query: %v", err)
	}
	return c.JSON(http.StatusOK, todo)
}
