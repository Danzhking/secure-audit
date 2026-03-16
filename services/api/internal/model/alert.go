package model

import "time"

type Alert struct {
	ID         int64      `json:"id"`
	RuleName   string     `json:"rule_name"`
	Severity   string     `json:"severity"`
	Message    string     `json:"message"`
	UserID     string     `json:"user_id"`
	IP         string     `json:"ip"`
	EventCount int        `json:"event_count"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}
