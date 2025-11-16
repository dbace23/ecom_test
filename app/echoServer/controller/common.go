package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func respondError(c echo.Context, code int, msg string, detail any) error {
	return c.JSON(code, echo.Map{
		"message": msg,
		"detail":  detail,
	})
}

func respondOK(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
		"data":    data,
	})
}
