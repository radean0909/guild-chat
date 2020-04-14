package models

import (
	"time"
)

type Conversation struct {
	ID        string     `json:"ID"`
	Sender    string     `json:"sender"`
	Recipient string     `json:"recipient"`
	Updated   time.Time  `json:"updated"`
	Messages  []*Message `json:"messages"`
}
