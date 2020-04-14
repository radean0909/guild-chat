package models

import "time"

type User struct {
	ID         string     `json:"id,omitempty"`
	Username   string     `json:"username,omitempty"`
	Email      string     `json:"email,omitempty"`
	ArchivedOn *time.Time `json:"archived_on,omitempty"`
}
