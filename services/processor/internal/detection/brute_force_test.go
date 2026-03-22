package detection

import (
	"errors"
	"testing"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
)

type mockCounter struct {
	failedByUser int
	failedByIP   int
	err          error
}

func (m *mockCounter) CountFailedLoginsByUser(userID string, windowMinutes int) (int, error) {
	return m.failedByUser, m.err
}

func (m *mockCounter) CountFailedLoginsByIP(ip string, windowMinutes int) (int, error) {
	return m.failedByIP, m.err
}

func TestBruteForceRule_Name(t *testing.T) {
	rule := NewBruteForceRule(&mockCounter{})
	if rule.Name() != "brute_force" {
		t.Errorf("expected 'brute_force', got '%s'", rule.Name())
	}
}

func TestBruteForceRule_SkipsIrrelevantEvents(t *testing.T) {
	rule := NewBruteForceRule(&mockCounter{})
	cases := []struct {
		name  string
		event model.Event
	}{
		{"login_success", model.Event{EventType: "login_success", UserID: "admin"}},
		{"file_access", model.Event{EventType: "file_access", UserID: "admin"}},
		{"empty_user", model.Event{EventType: "login_failed", UserID: ""}},
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

func TestBruteForceRule_NoAlertBelowThreshold(t *testing.T) {
	rule := NewBruteForceRule(&mockCounter{failedByUser: 3})
	alert, err := rule.Check(model.Event{EventType: "login_failed", UserID: "admin", IP: "1.2.3.4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert != nil {
		t.Error("expected no alert below threshold")
	}
}

func TestBruteForceRule_AlertAtThreshold(t *testing.T) {
	rule := NewBruteForceRule(&mockCounter{failedByUser: 5})
	alert, err := rule.Check(model.Event{EventType: "login_failed", UserID: "admin", IP: "1.2.3.4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert == nil {
		t.Fatal("expected alert at threshold")
	}
	if alert.RuleName != "brute_force" {
		t.Errorf("expected rule_name 'brute_force', got '%s'", alert.RuleName)
	}
	if alert.Severity != model.SeverityHigh {
		t.Errorf("expected severity 'high', got '%s'", alert.Severity)
	}
	if alert.UserID != "admin" {
		t.Errorf("expected user_id 'admin', got '%s'", alert.UserID)
	}
	if alert.EventCount != 5 {
		t.Errorf("expected event_count 5, got %d", alert.EventCount)
	}
}

func TestBruteForceRule_AlertAboveThreshold(t *testing.T) {
	rule := NewBruteForceRule(&mockCounter{failedByUser: 15})
	alert, err := rule.Check(model.Event{EventType: "login_failed", UserID: "victim", IP: "10.0.0.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alert == nil {
		t.Fatal("expected alert above threshold")
	}
	if alert.EventCount != 15 {
		t.Errorf("expected event_count 15, got %d", alert.EventCount)
	}
}

func TestBruteForceRule_RepoError(t *testing.T) {
	rule := NewBruteForceRule(&mockCounter{err: errors.New("db connection lost")})
	_, err := rule.Check(model.Event{EventType: "login_failed", UserID: "admin", IP: "1.2.3.4"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}
