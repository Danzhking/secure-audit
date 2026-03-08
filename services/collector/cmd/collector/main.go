package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Danzhking/secure-audit/services/collector/internal/config"
	"github.com/Danzhking/secure-audit/services/collector/internal/handler"
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

	r := gin.Default()

	r.POST("/events", eventHandler.CreateEvent)

	log.Println("Collector started on :8080")

	r.Run(":8080")
}
