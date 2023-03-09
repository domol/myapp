package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Handler struct {
	todoRepo Repo[Todo]
}

func (h Handler) list(c echo.Context) error {
	todos, err := h.todoRepo.list()
	if err != nil {
		return getServerErrorResponse(c, "Error connecting to the database.")
	}

	return c.JSON(http.StatusOK, todos)
}

func (h Handler) create(c echo.Context) error {
	var data Todo
	err := c.Bind(&data)
	if err != nil {
		return getBadRequestResponse(c, "Error parsing data.")
	}
	valid, err := govalidator.ValidateStruct(data)
	if !valid {
		return getBadRequestResponse(c, err.Error())
	}

	res, err := h.todoRepo.create(data)
	if err != nil {
		return getServerErrorResponse(c, "Database error.")
	}

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return getBadRequestResponse(c, "ID must be integer.")
	}

	err = h.todoRepo.delete(int64(id))
	if err != nil {
		return getServerErrorResponse(c, "DB error response.")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h Handler) retrieve(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return getBadRequestResponse(c, "ID must be integer.")
	}

	todo, err := h.todoRepo.get(int64(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return getServerErrorResponse(c, "Failed to fetch object.")
	}

	return c.JSON(http.StatusOK, todo)
}

func (h Handler) update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return getBadRequestResponse(c, "ID must be integer.")
	}

	var data Todo
	err = c.Bind(&data)
	if err != nil {
		return getBadRequestResponse(c, "Error parsing data.")
	}

	valid, err := govalidator.ValidateStruct(data)
	if !valid {
		return getBadRequestResponse(c, err.Error())
	}

	err = h.todoRepo.update(int64(id), data)
	if err != nil {
		return getServerErrorResponse(c, fmt.Sprintf("Couldn't update todo: %v", id))
	}

	todo, err := h.todoRepo.get(int64(id))
	if err != nil {
		return getServerErrorResponse(c, "Couldn't get updated data.")
	}
	return c.JSON(http.StatusOK, todo)
}
