package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/radean0909/guild-chat/api"
	"github.com/radean0909/guild-chat/api/models"
)


func GetMessageByID(svc *api.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		msg, err := svc.DB.GetMessage(id)

		return c.JSON(http.StatusOK, nil)
	}
}
