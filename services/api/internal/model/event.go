package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        int64           `json:"id"`
	Service   string          `json:"service"`
	EventType string          `json:"event_type"`
	Severity  string          `json:"severity"`
	UserID    string          `json:"user_id"`
	IP        string          `json:"ip"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
}
