package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Danzhking/secure-audit/services/api/internal/model"
)

type EventFilter struct {
	Service   string
	EventType string
	Severity  string
	UserID    string
	IP        string
	From      string
	To        string
	Page      int
	PageSize  int
}

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) List(f EventFilter) ([]model.Event, int, error) {
	where, args := buildEventWhere(f)

	countQuery := "SELECT COUNT(*) FROM security_events" + where
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.PageSize
	dataQuery := fmt.Sprintf(
		"SELECT id, service, event_type, severity, user_id, ip, COALESCE(metadata, '{}'), created_at FROM security_events%s ORDER BY created_at DESC LIMIT %d OFFSET %d",
		where, f.PageSize, offset,
	)

	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		if err := rows.Scan(&e.ID, &e.Service, &e.EventType, &e.Severity, &e.UserID, &e.IP, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}

	return events, total, rows.Err()
}

func (r *EventRepository) GetByID(id int64) (*model.Event, error) {
	query := "SELECT id, service, event_type, severity, user_id, ip, COALESCE(metadata, '{}'), created_at FROM security_events WHERE id = $1"

	var e model.Event
	err := r.db.QueryRow(query, id).Scan(&e.ID, &e.Service, &e.EventType, &e.Severity, &e.UserID, &e.IP, &e.Metadata, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func buildEventWhere(f EventFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	n := 1

	if f.Service != "" {
		conditions = append(conditions, fmt.Sprintf("service = $%d", n))
		args = append(args, f.Service)
		n++
	}
	if f.EventType != "" {
		conditions = append(conditions, fmt.Sprintf("event_type = $%d", n))
		args = append(args, f.EventType)
		n++
	}
	if f.Severity != "" {
		conditions = append(conditions, fmt.Sprintf("severity = $%d", n))
		args = append(args, f.Severity)
		n++
	}
	if f.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", n))
		args = append(args, f.UserID)
		n++
	}
	if f.IP != "" {
		conditions = append(conditions, fmt.Sprintf("ip = $%d", n))
		args = append(args, f.IP)
		n++
	}
	if f.From != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", n))
		args = append(args, f.From)
		n++
	}
	if f.To != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", n))
		args = append(args, f.To)
		n++
	}

	if len(conditions) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}
