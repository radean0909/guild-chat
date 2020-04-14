package models

import "time"

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	ArchivedOn time.Time `json:"archived_on"`
}
