package models

import "time"

type Message struct {
	ID        string    `json:"ID"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Date      time.Time `json:"date"`
}
