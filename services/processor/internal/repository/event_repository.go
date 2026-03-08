package repository

import (
	"database/sql"
	"log"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS security_events (
		id          BIGSERIAL PRIMARY KEY,
		service     VARCHAR(255) NOT NULL,
		event_type  VARCHAR(255) NOT NULL,
		user_id     VARCHAR(255) NOT NULL DEFAULT '',
		ip          VARCHAR(45)  NOT NULL DEFAULT '',
		created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_security_events_service    ON security_events (service);
	CREATE INDEX IF NOT EXISTS idx_security_events_event_type ON security_events (event_type);
	CREATE INDEX IF NOT EXISTS idx_security_events_user_id    ON security_events (user_id);
	CREATE INDEX IF NOT EXISTS idx_security_events_created_at ON security_events (created_at);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}

func (r *EventRepository) Save(event model.Event) error {
	query := `
		INSERT INTO security_events (service, event_type, user_id, ip)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRow(query,
		event.Service,
		event.EventType,
		event.UserID,
		event.IP,
	).Scan(&event.ID, &event.CreatedAt)
}
