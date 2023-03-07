package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type ErrorResponse struct {
	detail string
}

func getServerErrorResponse(c echo.Context, detail string) error {
	return c.JSON(http.StatusServiceUnavailable, ErrorResponse{detail: detail})
}

func getBadRequestResponse(c echo.Context, detail string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{detail: detail})
}
