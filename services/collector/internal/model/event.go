package model

type Event struct {
	Service   string                 `json:"service" binding:"required,max=255"`
	EventType string                 `json:"event_type" binding:"required,max=255"`
	Severity  string                 `json:"severity" binding:"omitempty,oneof=low medium high critical"`
	UserID    string                 `json:"user_id" binding:"max=255"`
	IP        string                 `json:"ip" binding:"omitempty,ip"`
	Metadata  map[string]interface{} `json:"metadata"`
}
