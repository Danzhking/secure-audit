package service

import (
	"encoding/json"

	"github.com/Danzhking/secure-audit/services/processor/internal/detection"
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
		var event model.Event

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			zap.L().Error("Failed to unmarshal message", zap.Error(err))
			msg.Nack(false, false)
			continue
		}

		if !event.Severity.IsValid() {
			event.Severity = model.SeverityLow
		}

		zap.L().Info("Processing event",
			zap.String("event_service", event.Service),
			zap.String("event_type", event.EventType),
			zap.String("severity", string(event.Severity)),
			zap.String("user_id", event.UserID),
			zap.String("ip", event.IP),
		)

		if err := s.repo.Save(event); err != nil {
			zap.L().Error("Failed to save event", zap.Error(err))
			msg.Nack(false, true)
			continue
		}

		s.engine.Analyze(event)

		msg.Ack(false)
		zap.L().Info("Event saved",
			zap.String("event_service", event.Service),
			zap.String("event_type", event.EventType),
			zap.String("severity", string(event.Severity)),
		)
	}
}
