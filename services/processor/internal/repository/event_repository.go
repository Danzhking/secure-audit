package repository

import (
	"database/sql"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Save(event model.Event) error {
	query := `
		INSERT INTO security_events (service, event_type, severity, user_id, ip, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	metadata := event.Metadata
	if len(metadata) == 0 {
		metadata = []byte("{}")
	}

	return r.db.QueryRow(query,
		event.Service,
		event.EventType,
		event.Severity,
		event.UserID,
		event.IP,
		metadata,
	).Scan(&event.ID, &event.CreatedAt)
}

func (r *EventRepository) CountFailedLoginsByUser(userID string, windowMinutes int) (int, error) {
	query := `
		SELECT COUNT(*) FROM security_events
		WHERE user_id = $1
		AND event_type = 'login_failed'
		AND created_at > NOW() - MAKE_INTERVAL(mins := $2)
	`

	var count int
	err := r.db.QueryRow(query, userID, windowMinutes).Scan(&count)
	return count, err
}

func (r *EventRepository) CountFailedLoginsByIP(ip string, windowMinutes int) (int, error) {
	query := `
		SELECT COUNT(DISTINCT user_id) FROM security_events
		WHERE ip = $1
		AND event_type = 'login_failed'
		AND created_at > NOW() - MAKE_INTERVAL(mins := $2)
	`

	var count int
	err := r.db.QueryRow(query, ip, windowMinutes).Scan(&count)
	return count, err
}
