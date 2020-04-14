package db

import (
	"time"

	"github.com/radean0909/guild-chat/api/models"
)

type Driver interface {
	GetMessage(id string) (*models.Message, error)
	CreateMessage(msg *models.Message) (*models.Message, error)
	ListMessages(recipient string, from, until time.Time, limit int) ([]*models.Message, error)
	GetConversation(sender, recipient string, from, until time.Time) (*models.Conversation, error)
	CreateConversation(sender, recipient string) (*models.Conversation, error)
	ListConversations(recipient string, from, until time.Time) ([]*models.Conversation, error)
	GetUser(id string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
}
