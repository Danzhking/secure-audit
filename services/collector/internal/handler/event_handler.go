package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Danzhking/secure-audit/services/collector/internal/metrics"
	"github.com/Danzhking/secure-audit/services/collector/internal/model"
	"github.com/Danzhking/secure-audit/services/collector/internal/service"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(s *service.EventService) *EventHandler {
	return &EventHandler{service: s}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {

	var event model.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		metrics.EventsReceived.WithLabelValues("rejected").Inc()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.ProcessEvent(event)

	if err != nil {
		metrics.EventsReceived.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process event",
		})
		return
	}

	metrics.EventsReceived.WithLabelValues("accepted").Inc()
	metrics.EventsPublished.Inc()
	c.JSON(http.StatusOK, gin.H{
		"status": "event accepted",
	})
}
