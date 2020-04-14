package models

import "time"

type Message struct {
	ID        string    `json:"ID"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
}
