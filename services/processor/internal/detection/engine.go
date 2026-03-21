package detection

import (
	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
	"go.uber.org/zap"
)

type Rule interface {
	Name() string
	Check(event model.Event) (*model.Alert, error)
}

type Engine struct {
	rules     []Rule
	alertRepo *repository.AlertRepository
}

func NewEngine(alertRepo *repository.AlertRepository, rules ...Rule) *Engine {
	names := make([]string, len(rules))
	for i, r := range rules {
		names[i] = r.Name()
	}
	zap.L().Info("Detection engine initialized",
		zap.Int("rule_count", len(rules)),
		zap.Strings("rules", names),
	)

	return &Engine{
		rules:     rules,
		alertRepo: alertRepo,
	}
}

func (e *Engine) Analyze(event model.Event) {
	for _, rule := range e.rules {
		alert, err := rule.Check(event)
		if err != nil {
			zap.L().Error("Detection rule error",
				zap.String("rule", rule.Name()),
				zap.Error(err),
			)
			continue
		}

		if alert == nil {
			continue
		}

		exists, err := e.alertRepo.HasActiveAlert(alert.RuleName, alert.UserID, alert.IP)
		if err != nil {
			zap.L().Error("Failed to check existing alert", zap.Error(err))
			continue
		}
		if exists {
			continue
		}

		id, err := e.alertRepo.Save(*alert)
		if err != nil {
			zap.L().Error("Failed to save alert", zap.Error(err))
			continue
		}

		zap.L().Warn("ALERT triggered",
			zap.Int64("alert_id", id),
			zap.String("rule", alert.RuleName),
			zap.String("severity", string(alert.Severity)),
			zap.String("message", alert.Message),
			zap.String("user_id", alert.UserID),
			zap.String("ip", alert.IP),
			zap.Int("event_count", alert.EventCount),
		)
	}
}
