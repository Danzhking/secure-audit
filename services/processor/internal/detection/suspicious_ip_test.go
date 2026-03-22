package detection

import (
	"errors"
	"testing"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
)

func TestSuspiciousIPRule_Name(t *testing.T) {
	rule := NewSuspiciousIPRule(&mockCounter{})
	if rule.Name() != "suspicious_ip" {
		t.Errorf("expected 'suspicious_ip', got '%s'", rule.Name())
	}
}

func TestSuspiciousIPRule_SkipsIrrelevantEvents(t *testing.T) {
	rule := NewSuspiciousIPRule(&mockCounter{})
	cases := []struct {
		name  string
		event model.Event
	}{
		{"login_success", model.Event{EventType: "login_success", IP: "1.2.3.4"}},
		{"empty_ip", model.Event{EventType: "login_failed", IP: ""}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			alert, err := rule.Check(tc.event)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if alert != nil {
				t.Errorf("expected nil alert")
			}
		})
	}
}

func TestSuspiciousIPRule_NoAlertBelowThreshold(t *testing.T) {
	rule := NewSuspiciousIPRule(&mockCounter{failedByIP: 2})
	alert, err := rule.Check(model.Event{EventType: "login_failed", UserID: "user1", IP: "10.0.0.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert != nil {
		t.Error("expected no alert below threshold")
	}
}

func TestSuspiciousIPRule_AlertAtThreshold(t *testing.T) {
	rule := NewSuspiciousIPRule(&mockCounter{failedByIP: 3})
	alert, err := rule.Check(model.Event{EventType: "login_failed", UserID: "user1", IP: "10.0.0.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert == nil {
		t.Fatal("expected alert at threshold")
	}
	if alert.RuleName != "suspicious_ip" {
		t.Errorf("expected rule_name 'suspicious_ip', got '%s'", alert.RuleName)
	}
	if alert.Severity != model.SeverityCritical {
		t.Errorf("expected severity 'critical', got '%s'", alert.Severity)
	}
	if alert.UserID != "" {
		t.Errorf("expected empty user_id, got '%s'", alert.UserID)
	}
	if alert.IP != "10.0.0.1" {
		t.Errorf("expected ip '10.0.0.1', got '%s'", alert.IP)
	}
}

func TestSuspiciousIPRule_RepoError(t *testing.T) {
	rule := NewSuspiciousIPRule(&mockCounter{err: errors.New("timeout")})
	_, err := rule.Check(model.Event{EventType: "login_failed", UserID: "user1", IP: "10.0.0.1"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}
