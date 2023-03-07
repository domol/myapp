package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type ErrorResponse struct {
	Detail string `json:"detail" xml:"detail"`
}

func getServerErrorResponse(c echo.Context, detail string) error {
	return c.JSON(http.StatusServiceUnavailable, ErrorResponse{Detail: detail})
}

func getBadRequestResponse(c echo.Context, detail string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{Detail: detail})
}
