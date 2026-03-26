package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		zap.L().Fatal("Не удалось создать издателя RabbitMQ", zap.Error(err))
	}

	eventService := service.NewEventService(publisher)
	eventHandler := handler.NewEventHandler(eventService)

	rateLimiter := middleware.NewRateLimiter(10, 20)

	// Prometheus metrics on separate HTTP port (no TLS, no auth)
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		zap.L().Info("Сервер метрик запущен", zap.String("port", ":9090"))
		if err := http.ListenAndServe(":9090", mux); err != nil {
			zap.L().Error("Сервер метрик завершился с ошибкой", zap.Error(err))
		}
	}()

	r := gin.Default()

	r.POST("/events",
		rateLimiter.Middleware(),
		middleware.APIKeyAuth(cfg.APIKeys),
		middleware.HMACVerify(cfg.HMACSecret),
		eventHandler.CreateEvent,
	)

	if cfg.TLSEnabled() {
		zap.L().Info("Collector запущен с TLS",
			zap.String("port", cfg.TLSPort),
			zap.String("cert", cfg.TLSCert),
		)
		if err := r.RunTLS(cfg.TLSPort, cfg.TLSCert, cfg.TLSKey); err != nil {
			zap.L().Fatal("Не удалось запустить TLS-сервер", zap.Error(err))
		}
	} else {
		zap.L().Info("Collector запущен без TLS", zap.String("port", cfg.Port))
		r.Run(cfg.Port)
	}
}
