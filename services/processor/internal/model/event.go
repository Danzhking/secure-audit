package model

import (
	"encoding/json"
	"time"
)

type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

func (s Severity) IsValid() bool {
	switch s {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	}
	return false
}

type Event struct {
	ID        int64           `json:"id" db:"id"`
	Service   string          `json:"service" db:"service"`
	EventType string          `json:"event_type" db:"event_type"`
	Severity  Severity        `json:"severity" db:"severity"`
	UserID    string          `json:"user_id" db:"user_id"`
	IP        string          `json:"ip" db:"ip"`
	Metadata  json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
