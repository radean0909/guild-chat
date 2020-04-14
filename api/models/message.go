package models

import "time"

type Message struct {
	ID        string     `json:"id,omitempty"`
	Sender    string     `json:"sender,omitempty"`
	Recipient string     `json:"recipient,omitempty"`
	Content   string     `json:"content,omitempty"`
	Date      *time.Time `json:"date,omitempty"`
}
