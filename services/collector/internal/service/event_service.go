package service

import (
	"github.com/Danzhking/secure-audit/services/collector/internal/model"
	"github.com/Danzhking/secure-audit/services/collector/internal/queue"
	"go.uber.org/zap"
)

type EventService struct {
	publisher *queue.Publisher
}

func NewEventService(p *queue.Publisher) *EventService {
	return &EventService{
		publisher: p,
	}
}

func (s *EventService) ProcessEvent(event model.Event) error {
	zap.L().Info("Обработка события",
		zap.String("event_type", event.EventType),
		zap.String("event_service", event.Service),
		zap.String("user_id", event.UserID),
		zap.String("ip", event.IP),
	)
	return s.publisher.Publish(event)
}
