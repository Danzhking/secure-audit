package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Danzhking/secure-audit/services/api/internal/config"
	"github.com/Danzhking/secure-audit/services/api/internal/handler"
	"github.com/Danzhking/secure-audit/services/api/internal/repository"
)

func main() {
	cfg := config.Load()

	db := repository.ConnectPostgres(cfg.PostgresURL)
	defer db.Close()

	eventRepo := repository.NewEventRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	eventHandler := handler.NewEventHandler(eventRepo)
	alertHandler := handler.NewAlertHandler(alertRepo)
	statsHandler := handler.NewStatsHandler(statsRepo)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/events", eventHandler.List)
		api.GET("/events/:id", eventHandler.GetByID)

		api.GET("/alerts", alertHandler.List)
		api.PATCH("/alerts/:id", alertHandler.UpdateStatus)

		api.GET("/stats", statsHandler.GetStats)
	}

	log.Printf("API service started on %s", cfg.Port)
	r.Run(cfg.Port)
}
