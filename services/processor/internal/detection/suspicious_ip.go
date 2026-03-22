package detection

import (
	"fmt"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
)

type SuspiciousIPRule struct {
	counter       EventCounter
	Threshold     int
	WindowMinutes int
}

func NewSuspiciousIPRule(counter EventCounter) *SuspiciousIPRule {
	return &SuspiciousIPRule{
		counter:       counter,
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

	distinctUsers, err := r.counter.CountFailedLoginsByIP(event.IP, r.WindowMinutes)
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
