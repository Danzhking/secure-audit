package detection

import (
	"fmt"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
)

// BruteForceRule detects multiple failed login attempts from the same user.
// Triggers when a user has >= Threshold failed logins within WindowMinutes.
type BruteForceRule struct {
	eventRepo     *repository.EventRepository
	Threshold     int
	WindowMinutes int
}

func NewBruteForceRule(eventRepo *repository.EventRepository) *BruteForceRule {
	return &BruteForceRule{
		eventRepo:     eventRepo,
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

	count, err := r.eventRepo.CountFailedLoginsByUser(event.UserID, r.WindowMinutes)
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
