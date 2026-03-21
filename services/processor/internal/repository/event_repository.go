package repository

import (
	"database/sql"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"go.uber.org/zap"
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
		id          BIGSERIAL    PRIMARY KEY,
		service     VARCHAR(255) NOT NULL,
		event_type  VARCHAR(255) NOT NULL,
		severity    VARCHAR(20)  NOT NULL DEFAULT 'low',
		user_id     VARCHAR(255) NOT NULL DEFAULT '',
		ip          VARCHAR(45)  NOT NULL DEFAULT '',
		metadata    JSONB        DEFAULT '{}',
		created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);

	ALTER TABLE security_events ADD COLUMN IF NOT EXISTS severity VARCHAR(20) NOT NULL DEFAULT 'low';
	ALTER TABLE security_events ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

	CREATE INDEX IF NOT EXISTS idx_security_events_service    ON security_events (service);
	CREATE INDEX IF NOT EXISTS idx_security_events_event_type ON security_events (event_type);
	CREATE INDEX IF NOT EXISTS idx_security_events_severity   ON security_events (severity);
	CREATE INDEX IF NOT EXISTS idx_security_events_user_id    ON security_events (user_id);
	CREATE INDEX IF NOT EXISTS idx_security_events_ip         ON security_events (ip);
	CREATE INDEX IF NOT EXISTS idx_security_events_created_at ON security_events (created_at);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	zap.L().Info("Database migration completed")
	return nil
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
