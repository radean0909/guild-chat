package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/radean0909/guild-chat/api/internal/constants"
	"github.com/radean0909/guild-chat/api/internal/db"
	"github.com/radean0909/guild-chat/api/models"
)

type ConversationHandler struct {
	DB db.Driver
}

// GetConversation - returns all messages sent from a person to another person.
func (h *ConversationHandler) GetConversation(c echo.Context) error {
	sender := c.Param("from")
	recipient := c.Param("to")

	// additional query params (to satisfy challenege requirements)
	startParam := c.QueryParam("start")
	untilParam := c.QueryParam("until")
	limitParam := c.QueryParam("limit")

	var start, until time.Time
	var err error
	limit := 100 // set a default limit to 100
	if startParam != "" {
		start, err = time.Parse("2006-01-02", startParam)
		if err != nil {
			return handleError(c, err)
		}
	} else {
		start = time.Now().AddDate(0, 0, -30)
	}

	if untilParam != "" {
		until, err = time.Parse("2006-01-02", untilParam)
		if err != nil {
			return handleError(c, err)
		}
	} else {
		until = time.Now()
	}

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return handleError(c, err)
		}
	}

	convo, err := h.DB.GetConversation(sender, recipient, start, until)

	if err != nil {
		return handleError(c, err)
	}

	// Now filter the messages
	messages := []*models.Message{}
	for _, msg := range convo.Messages {
		if len(messages) >= limit {
			break
		}
		if msg.Recipient == recipient {
			messages = append(messages, msg)
		}
	}

	if len(messages) > 0 {
		return c.JSON(http.StatusOK, messages)
	}

	// this shouldn't normally happen, but if there are orphaned conversations, we don't want to return empty messages
	return handleError(c, constants.ErrNotFound)

}

// ListConversations - lists all messages sent to a particular person
func (h *ConversationHandler) ListConversations(c echo.Context) error {
	recipient := c.Param("to")

	// additional query params (to satisfy challenege requirements)
	startParam := c.QueryParam("start")
	untilParam := c.QueryParam("until")
	limitParam := c.QueryParam("limit")

	var start, until time.Time
	var err error
	limit := 100 // set a default limit to 100
	if startParam != "" {
		start, err = time.Parse("2006-01-02", startParam)
		if err != nil {
			return handleError(c, err)
		}
	} else {
		start = time.Now().AddDate(0, 0, -30)
	}

	if untilParam != "" {
		until, err = time.Parse("2006-01-02", untilParam)
		if err != nil {
			return handleError(c, err)
		}
	} else {
		until = time.Now()
	}

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			return handleError(c, err)
		}
	}

	convos, err := h.DB.ListConversations(recipient, start, until)

	if err != nil {
		return handleError(c, err)
	}

	// Now filter the messages
	messages := []*models.Message{}
	for _, convo := range convos {
		for _, msg := range convo.Messages {
			if len(messages) >= limit {
				break
			}
			if msg.Recipient == recipient {
				messages = append(messages, msg)
			}
		}
	}

	if len(messages) > 0 {
		return c.JSON(http.StatusOK, messages)
	}

	// this shouldn't normally happen, but if there are orphaned conversations, we don't want to return empty messages
	return handleError(c, constants.ErrNotFound)
}
