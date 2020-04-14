package db

import (
	"database/sql"

	"github.com/radean0909/guild-chat/api/models"
)

type Driver interface {
	GetMessage(id string) (error, *models.Message)
	
}
