package service

import (
	"log"

	"github.com/Danzhking/secure-audit/services/collector/internal/model"
	"github.com/Danzhking/secure-audit/services/collector/internal/queue"
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

	log.Println("Processing event:", event.EventType)
	return s.publisher.Publish(event)
}
