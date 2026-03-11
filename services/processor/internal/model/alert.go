package model

import "time"

type AlertStatus string

const (
	AlertStatusNew          AlertStatus = "new"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved     AlertStatus = "resolved"
)

type Alert struct {
	ID         int64       `json:"id" db:"id"`
	RuleName   string      `json:"rule_name" db:"rule_name"`
	Severity   Severity    `json:"severity" db:"severity"`
	Message    string      `json:"message" db:"message"`
	UserID     string      `json:"user_id" db:"user_id"`
	IP         string      `json:"ip" db:"ip"`
	EventCount int         `json:"event_count" db:"event_count"`
	Status     AlertStatus `json:"status" db:"status"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	ResolvedAt *time.Time  `json:"resolved_at,omitempty" db:"resolved_at"`
}
