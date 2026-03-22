package detection

import (
	"fmt"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
)

type EventCounter interface {
	CountFailedLoginsByUser(userID string, windowMinutes int) (int, error)
	CountFailedLoginsByIP(ip string, windowMinutes int) (int, error)
}

type BruteForceRule struct {
	counter       EventCounter
	Threshold     int
	WindowMinutes int
}

func NewBruteForceRule(counter EventCounter) *BruteForceRule {
	return &BruteForceRule{
		counter:       counter,
		Threshold:     5,
		WindowMinutes: 10,
	}
}

func (r *BruteForceRule) Name() string {
	return "brute_force"
}

func (r *BruteForceRule) Check(event model.Event) (*model.Alert, error) {
	if event.EventType != "login_failed" || event.UserID == "" {
		return nil, nil
	}

	count, err := r.counter.CountFailedLoginsByUser(event.UserID, r.WindowMinutes)
	if err != nil {
		return nil, fmt.Errorf("count failed logins: %w", err)
	}

	if count < r.Threshold {
		return nil, nil
	}

	return &model.Alert{
		RuleName:   r.Name(),
		Severity:   model.SeverityHigh,
		Message:    fmt.Sprintf("Brute force detected: user '%s' has %d failed login attempts in %d minutes", event.UserID, count, r.WindowMinutes),
		UserID:     event.UserID,
		IP:         event.IP,
		EventCount: count,
	}, nil
}
