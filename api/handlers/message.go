package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type Message struct {
	ID        string    `json:"ID"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Date      time.Time `json:"date"`
}

func GetMessageByID(c echo.Context) error {
	id := c.Param("id")

	msg := &Message{}

	return c.JSON(http.StatusOK, nil)
}
