package detection

import (
	"fmt"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
)

// SuspiciousIPRule detects credential scanning — a single IP targeting
// multiple distinct user accounts with failed logins.
// Triggers when an IP has failed logins against >= Threshold different users within WindowMinutes.
type SuspiciousIPRule struct {
	eventRepo     *repository.EventRepository
	Threshold     int
	WindowMinutes int
}

func NewSuspiciousIPRule(eventRepo *repository.EventRepository) *SuspiciousIPRule {
	return &SuspiciousIPRule{
		eventRepo:     eventRepo,
		Threshold:     3,
		WindowMinutes: 5,
	}
}

func (r *SuspiciousIPRule) Name() string {
	return "suspicious_ip"
}

func (r *SuspiciousIPRule) Check(event model.Event) (*model.Alert, error) {
	if event.EventType != "login_failed" || event.IP == "" {
		return nil, nil
	}

	distinctUsers, err := r.eventRepo.CountFailedLoginsByIP(event.IP, r.WindowMinutes)
	if err != nil {
		return nil, fmt.Errorf("count failed logins by IP: %w", err)
	}

	if distinctUsers < r.Threshold {
		return nil, nil
	}

	return &model.Alert{
		RuleName:   r.Name(),
		Severity:   model.SeverityCritical,
		Message:    fmt.Sprintf("Credential scanning detected: IP '%s' targeted %d distinct users in %d minutes", event.IP, distinctUsers, r.WindowMinutes),
		UserID:     "",
		IP:         event.IP,
		EventCount: distinctUsers,
	}, nil
}
