package model

import "time"

type Event struct {
	ID        int64     `json:"id" db:"id"`
	Service   string    `json:"service" db:"service"`
	EventType string    `json:"event_type" db:"event_type"`
	UserID    string    `json:"user_id" db:"user_id"`
	IP        string    `json:"ip" db:"ip"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
