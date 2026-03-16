package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Danzhking/secure-audit/services/api/internal/model"
)

type AlertFilter struct {
	RuleName string
	Severity string
	Status   string
	Page     int
	PageSize int
}

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) List(f AlertFilter) ([]model.Alert, int, error) {
	where, args := buildAlertWhere(f)

	countQuery := "SELECT COUNT(*) FROM alerts" + where
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.PageSize
	dataQuery := fmt.Sprintf(
		"SELECT id, rule_name, severity, message, user_id, ip, event_count, status, created_at, resolved_at FROM alerts%s ORDER BY created_at DESC LIMIT %d OFFSET %d",
		where, f.PageSize, offset,
	)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var alerts []model.Alert
	for rows.Next() {
		var a model.Alert
		if err := rows.Scan(&a.ID, &a.RuleName, &a.Severity, &a.Message, &a.UserID, &a.IP, &a.EventCount, &a.Status, &a.CreatedAt, &a.ResolvedAt); err != nil {
			return nil, 0, err
		}
		alerts = append(alerts, a)
	}

	return alerts, total, rows.Err()
}

func (r *AlertRepository) UpdateStatus(id int64, status string) error {
	var query string
	if status == "resolved" {
		query = "UPDATE alerts SET status = $1, resolved_at = $2 WHERE id = $3"
		_, err := r.db.Exec(query, status, time.Now(), id)
		return err
	}

	query = "UPDATE alerts SET status = $1 WHERE id = $2"
	_, err := r.db.Exec(query, status, id)
	return err
}

func buildAlertWhere(f AlertFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	n := 1

	if f.RuleName != "" {
		conditions = append(conditions, fmt.Sprintf("rule_name = $%d", n))
		args = append(args, f.RuleName)
		n++
	}
	if f.Severity != "" {
		conditions = append(conditions, fmt.Sprintf("severity = $%d", n))
		args = append(args, f.Severity)
		n++
	}
	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", n))
		args = append(args, f.Status)
		n++
	}

	if len(conditions) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}
