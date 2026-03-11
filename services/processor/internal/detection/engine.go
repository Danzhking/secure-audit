package detection

import (
	"log"

	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
)

type Rule interface {
	Name() string
	Check(event model.Event) (*model.Alert, error)
}

type Engine struct {
	rules    []Rule
	alertRepo *repository.AlertRepository
}

func NewEngine(alertRepo *repository.AlertRepository, rules ...Rule) *Engine {
	names := make([]string, len(rules))
	for i, r := range rules {
		names[i] = r.Name()
	}
	log.Printf("Detection engine initialized with %d rules: %v", len(rules), names)

	return &Engine{
		rules:    rules,
		alertRepo: alertRepo,
	}
}

func (e *Engine) Analyze(event model.Event) {
	for _, rule := range e.rules {
		alert, err := rule.Check(event)
		if err != nil {
			log.Printf("Detection rule '%s' error: %v", rule.Name(), err)
			continue
		}

		if alert == nil {
			continue
		}

		exists, err := e.alertRepo.HasActiveAlert(alert.RuleName, alert.UserID, alert.IP)
		if err != nil {
			log.Printf("Failed to check existing alert: %v", err)
			continue
		}
		if exists {
			continue
		}

		id, err := e.alertRepo.Save(*alert)
		if err != nil {
			log.Printf("Failed to save alert: %v", err)
			continue
		}

		log.Printf("ALERT #%d [%s] %s: %s", id, alert.Severity, alert.RuleName, alert.Message)
	}
}
