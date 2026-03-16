package repository

import "database/sql"

type Stats struct {
	TotalEvents  int              `json:"total_events"`
	TotalAlerts  int              `json:"total_alerts"`
	ActiveAlerts int              `json:"active_alerts"`
	UniqueUsers  int              `json:"unique_users"`
	UniqueIPs    int              `json:"unique_ips"`
	ByEventType  []CountEntry     `json:"by_event_type"`
	BySeverity   []CountEntry     `json:"by_severity"`
	TopIPs       []CountEntry     `json:"top_ips"`
	TopUsers     []CountEntry     `json:"top_users"`
}

type CountEntry struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetStats() (*Stats, error) {
	s := &Stats{}

	r.db.QueryRow("SELECT COUNT(*) FROM security_events").Scan(&s.TotalEvents)
	r.db.QueryRow("SELECT COUNT(*) FROM alerts").Scan(&s.TotalAlerts)
	r.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE status = 'new'").Scan(&s.ActiveAlerts)
	r.db.QueryRow("SELECT COUNT(DISTINCT user_id) FROM security_events WHERE user_id != ''").Scan(&s.UniqueUsers)
	r.db.QueryRow("SELECT COUNT(DISTINCT ip) FROM security_events WHERE ip != ''").Scan(&s.UniqueIPs)

	var err error

	s.ByEventType, err = r.queryGroup("SELECT event_type, COUNT(*) FROM security_events GROUP BY event_type ORDER BY COUNT(*) DESC LIMIT 10")
	if err != nil {
		return nil, err
	}

	s.BySeverity, err = r.queryGroup("SELECT severity, COUNT(*) FROM security_events GROUP BY severity ORDER BY COUNT(*) DESC")
	if err != nil {
		return nil, err
	}

	s.TopIPs, err = r.queryGroup("SELECT ip, COUNT(*) FROM security_events WHERE ip != '' GROUP BY ip ORDER BY COUNT(*) DESC LIMIT 10")
	if err != nil {
		return nil, err
	}

	s.TopUsers, err = r.queryGroup("SELECT user_id, COUNT(*) FROM security_events WHERE user_id != '' GROUP BY user_id ORDER BY COUNT(*) DESC LIMIT 10")
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *StatsRepository) queryGroup(query string) ([]CountEntry, error) {
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []CountEntry
	for rows.Next() {
		var e CountEntry
		if err := rows.Scan(&e.Name, &e.Count); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
