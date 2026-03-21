package repository

import (
	"database/sql"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"go.uber.org/zap"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS alerts (
		id          BIGSERIAL    PRIMARY KEY,
		rule_name   VARCHAR(255) NOT NULL,
		severity    VARCHAR(20)  NOT NULL,
		message     TEXT         NOT NULL,
		user_id     VARCHAR(255) DEFAULT '',
		ip          VARCHAR(45)  DEFAULT '',
		event_count INT          NOT NULL DEFAULT 0,
		status      VARCHAR(20)  NOT NULL DEFAULT 'new',
		created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		resolved_at TIMESTAMPTZ
	);

	CREATE INDEX IF NOT EXISTS idx_alerts_rule_name  ON alerts (rule_name);
	CREATE INDEX IF NOT EXISTS idx_alerts_severity   ON alerts (severity);
	CREATE INDEX IF NOT EXISTS idx_alerts_status     ON alerts (status);
	CREATE INDEX IF NOT EXISTS idx_alerts_user_id    ON alerts (user_id);
	CREATE INDEX IF NOT EXISTS idx_alerts_ip         ON alerts (ip);
	CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts (created_at);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	zap.L().Info("Alerts migration completed")
	return nil
}

func (r *AlertRepository) Save(alert model.Alert) (int64, error) {
	query := `
		INSERT INTO alerts (rule_name, severity, message, user_id, ip, event_count)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(query,
		alert.RuleName,
		alert.Severity,
		alert.Message,
		alert.UserID,
		alert.IP,
		alert.EventCount,
	).Scan(&id)

	return id, err
}

func (r *AlertRepository) HasActiveAlert(ruleName, userID, ip string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM alerts
			WHERE rule_name = $1
			AND ($2 = '' OR user_id = $2)
			AND ($3 = '' OR ip = $3)
			AND status = 'new'
			AND created_at > NOW() - INTERVAL '30 minutes'
		)
	`

	var exists bool
	err := r.db.QueryRow(query, ruleName, userID, ip).Scan(&exists)
	return exists, err
}
