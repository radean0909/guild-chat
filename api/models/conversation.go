package models

import (
	"time"
)

type Conversation struct {
	ID        string     `json:"id,omitempty"`
	Sender    string     `json:"sender,omitempty"`
	Recipient string     `json:"recipient,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
	Messages  []*Message `json:"messages,omitempty"`
}
