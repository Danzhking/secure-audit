package repository

import (
	"database/sql"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{db: db}
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
