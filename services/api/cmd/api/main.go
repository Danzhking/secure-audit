package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Danzhking/secure-audit/services/api/internal/config"
	"github.com/Danzhking/secure-audit/services/api/internal/handler"
	"github.com/Danzhking/secure-audit/services/api/internal/logger"
	"github.com/Danzhking/secure-audit/services/api/internal/middleware"
	"github.com/Danzhking/secure-audit/services/api/internal/repository"
)

func main() {
	logger.Init("api")
	defer zap.L().Sync()

	cfg := config.Load()

	db := repository.ConnectPostgres(cfg.PostgresURL)
	defer db.Close()

	eventRepo := repository.NewEventRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	eventHandler := handler.NewEventHandler(eventRepo)
	alertHandler := handler.NewAlertHandler(alertRepo)
	statsHandler := handler.NewStatsHandler(statsRepo)
	authHandler := handler.NewAuthHandler(cfg.JWTSecret)

	r := gin.Default()

	r.POST("/auth/login", authHandler.Login)

	api := r.Group("/api")
	api.Use(middleware.JWTAuth(cfg.JWTSecret))
	api.Use(middleware.AuditLog())
	{
		api.GET("/events", eventHandler.List)
		api.GET("/events/:id", eventHandler.GetByID)

		api.GET("/alerts", alertHandler.List)
		api.PATCH("/alerts/:id", alertHandler.UpdateStatus)

		api.GET("/stats", statsHandler.GetStats)
	}

	zap.L().Info("Сервис API запущен", zap.String("port", cfg.Port))
	r.Run(cfg.Port)
}
