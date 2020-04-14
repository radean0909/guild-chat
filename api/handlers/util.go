package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/radean0909/guild-chat/api/internal/constants"
)

func handleError(c echo.Context, err error) error {
	if errors.Is(err, constants.ErrBadRequest) {
		return c.JSON(http.StatusBadRequest, err)
	}
	if errors.Is(err, constants.ErrNotFound) {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusInternalServerError, err)
}
