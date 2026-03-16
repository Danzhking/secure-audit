package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Danzhking/secure-audit/services/collector/internal/config"
	"github.com/Danzhking/secure-audit/services/collector/internal/handler"
	"github.com/Danzhking/secure-audit/services/collector/internal/middleware"
	"github.com/Danzhking/secure-audit/services/collector/internal/queue"
	"github.com/Danzhking/secure-audit/services/collector/internal/service"
)

func main() {
	cfg := config.Load()

	conn := queue.ConnectRabbitMQ(cfg.RabbitURL)
	defer conn.Close()

	publisher, err := queue.NewPublisher(conn)
	if err != nil {
		log.Fatal(err)
	}

	eventService := service.NewEventService(publisher)
	eventHandler := handler.NewEventHandler(eventService)

	rateLimiter := middleware.NewRateLimiter(10, 20) // 10 req/s, burst 20

	r := gin.Default()

	r.POST("/events",
		rateLimiter.Middleware(),
		middleware.APIKeyAuth(cfg.APIKeys),
		middleware.HMACVerify(cfg.HMACSecret),
		eventHandler.CreateEvent,
	)

	log.Printf("Collector started on %s", cfg.Port)
	r.Run(cfg.Port)
}
