package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/Danzhking/secure-audit/services/processor/internal/config"
	"github.com/Danzhking/secure-audit/services/processor/internal/detection"
	"github.com/Danzhking/secure-audit/services/processor/internal/logger"
	"github.com/Danzhking/secure-audit/services/processor/internal/queue"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
	"github.com/Danzhking/secure-audit/services/processor/internal/service"
)

func main() {
	logger.Init("processor")
	defer zap.L().Sync()

	cfg := config.Load()

	db := repository.ConnectPostgres(cfg.PostgresURL)
	defer db.Close()

	eventRepo := repository.NewEventRepository(db)
	if err := eventRepo.Migrate(); err != nil {
		zap.L().Fatal("Ошибка миграции таблицы событий", zap.Error(err))
	}

	alertRepo := repository.NewAlertRepository(db)
	if err := alertRepo.Migrate(); err != nil {
		zap.L().Fatal("Ошибка миграции таблицы оповещений", zap.Error(err))
	}

	engine := detection.NewEngine(
		alertRepo,
		detection.NewBruteForceRule(eventRepo),
		detection.NewSuspiciousIPRule(eventRepo),
	)

	conn := queue.ConnectRabbitMQ(cfg.RabbitURL)
	defer conn.Close()

	consumer, err := queue.NewConsumer(conn)
	if err != nil {
		zap.L().Fatal("Не удалось создать потребителя RabbitMQ", zap.Error(err))
	}
	defer consumer.Close()

	msgs, err := consumer.Consume()
	if err != nil {
		zap.L().Fatal("Не удалось начать потребление сообщений", zap.Error(err))
	}

	eventService := service.NewEventService(eventRepo, engine)

	// Prometheus metrics
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		zap.L().Info("Сервер метрик запущен", zap.String("port", ":9091"))
		if err := http.ListenAndServe(":9091", mux); err != nil {
			zap.L().Error("Сервер метрик завершился с ошибкой", zap.Error(err))
		}
	}()

	go eventService.ProcessMessages(msgs)

	zap.L().Info("Processor запущен, ожидание событий")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	zap.L().Info("Завершение работы Processor")
}
