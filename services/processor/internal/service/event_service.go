package service

import (
	"encoding/json"
	"log"

	"github.com/Danzhking/secure-audit/services/processor/internal/detection"
	"github.com/Danzhking/secure-audit/services/processor/internal/model"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
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
			log.Printf("Failed to unmarshal message: %v", err)
			msg.Nack(false, false)
			continue
		}

		if !event.Severity.IsValid() {
			event.Severity = model.SeverityLow
		}

		log.Printf("Processing event: service=%s type=%s severity=%s user=%s ip=%s",
			event.Service, event.EventType, event.Severity, event.UserID, event.IP)

		if err := s.repo.Save(event); err != nil {
			log.Printf("Failed to save event: %v", err)
			msg.Nack(false, true)
			continue
		}

		s.engine.Analyze(event)

		msg.Ack(false)
		log.Printf("Event saved: %s/%s [%s]", event.Service, event.EventType, event.Severity)
	}
}
