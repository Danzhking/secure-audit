package service

import (
	"encoding/json"
	"time"

	"github.com/Danzhking/secure-audit/services/processor/internal/detection"
	"github.com/Danzhking/secure-audit/services/processor/internal/metrics"
	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type EventService struct {
	repo   *repository.EventRepository
	engine *detection.Engine
}

func NewEventService(repo *repository.EventRepository, engine *detection.Engine) *EventService {
	return &EventService{
		repo:   repo,
		engine: engine,
	}
}

func (s *EventService) ProcessMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		start := time.Now()

		var event model.Event

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			zap.L().Error("Не удалось разобрать сообщение", zap.Error(err))
			metrics.EventsProcessed.WithLabelValues("unmarshal_error").Inc()
			msg.Nack(false, false)
			continue
		}

		if !event.Severity.IsValid() {
			event.Severity = model.SeverityLow
		}

		zap.L().Info("Обработка события",
			zap.String("event_service", event.Service),
			zap.String("event_type", event.EventType),
			zap.String("severity", string(event.Severity)),
			zap.String("user_id", event.UserID),
			zap.String("ip", event.IP),
		)

		if err := s.repo.Save(event); err != nil {
			zap.L().Error("Не удалось сохранить событие", zap.Error(err))
			metrics.EventsProcessed.WithLabelValues("save_error").Inc()
			msg.Nack(false, true)
			continue
		}

		s.engine.Analyze(event)

		msg.Ack(false)
		metrics.EventsProcessed.WithLabelValues("success").Inc()
		metrics.ProcessingDuration.Observe(time.Since(start).Seconds())

		zap.L().Info("Событие сохранено",
			zap.String("event_service", event.Service),
			zap.String("event_type", event.EventType),
			zap.String("severity", string(event.Severity)),
		)
	}
}
