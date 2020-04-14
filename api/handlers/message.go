package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/radean0909/guild-chat/api/internal/db"
	"github.com/radean0909/guild-chat/api/models"
)

// MessageHandler - the message handler. In a full-fledged app, this might dial into a gRPC service, for instance
type MessageHandler struct {
	DB db.Driver
}

// GetMessageByID - GET: retrieve a single message by message ID
func (h *MessageHandler) GetMessageByID(c echo.Context) error {
	id := c.Param("id")

	msg, err := h.DB.GetMessage(id)

	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, msg)

}

// PostMessage - POST: create a new message
func (h *MessageHandler) PostMessage(c echo.Context) error {
	msg := &models.Message{}

	if err := c.Bind(msg); err != nil {
		return handleError(c, err)
	}

	msg, err := h.DB.CreateMessage(msg)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, msg)

}
