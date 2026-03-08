package model

type Event struct {
	Service   string                 `json:"service" binding:"required"`
	EventType string                 `json:"event_type" binding:"required"`
	UserID    string                 `json:"user_id"`
	IP        string                 `json:"ip"`
	Metadata  map[string]interface{} `json:"metadata"`
}
