package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Danzhking/secure-audit/services/collector/internal/config"
	"github.com/Danzhking/secure-audit/services/collector/internal/handler"
	"github.com/Danzhking/secure-audit/services/collector/internal/logger"
	"github.com/Danzhking/secure-audit/services/collector/internal/middleware"
	"github.com/Danzhking/secure-audit/services/collector/internal/queue"
	"github.com/Danzhking/secure-audit/services/collector/internal/service"
)

func main() {
	logger.Init("collector")
	defer zap.L().Sync()

	cfg := config.Load()

	conn := queue.ConnectRabbitMQ(cfg.RabbitURL)
	defer conn.Close()

	publisher, err := queue.NewPublisher(conn)
	if err != nil {
		zap.L().Fatal("Failed to create publisher", zap.Error(err))
	}

	eventService := service.NewEventService(publisher)
	eventHandler := handler.NewEventHandler(eventService)

	rateLimiter := middleware.NewRateLimiter(10, 20)

	r := gin.Default()

	r.POST("/events",
		rateLimiter.Middleware(),
		middleware.APIKeyAuth(cfg.APIKeys),
		middleware.HMACVerify(cfg.HMACSecret),
		eventHandler.CreateEvent,
	)

	if cfg.TLSEnabled() {
		zap.L().Info("Collector started with TLS",
			zap.String("port", cfg.TLSPort),
			zap.String("cert", cfg.TLSCert),
		)
		if err := r.RunTLS(cfg.TLSPort, cfg.TLSCert, cfg.TLSKey); err != nil {
			zap.L().Fatal("Failed to start TLS server", zap.Error(err))
		}
	} else {
		zap.L().Info("Collector started (no TLS)", zap.String("port", cfg.Port))
		r.Run(cfg.Port)
	}
}
