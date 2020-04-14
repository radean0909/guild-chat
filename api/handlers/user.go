package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/radean0909/guild-chat/api/internal/db"
	"github.com/radean0909/guild-chat/api/models"
)

type UserHandler struct {
	DB db.Driver
}

func (h *UserHandler) PostUser(c echo.Context) error {
	user := &models.User{}

	if err := c.Bind(user); err != nil {
		return handleError(c, err)
	}

	user, err := h.DB.CreateUser(user)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")

	user, err := h.DB.GetUser(id)

	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUserbyID(c echo.Context) error {
	id := c.Param("id")

	err := h.DB.DeleteUser(id)

	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}
