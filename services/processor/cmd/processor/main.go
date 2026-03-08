package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Danzhking/secure-audit/services/processor/internal/config"
	"github.com/Danzhking/secure-audit/services/processor/internal/queue"
	"github.com/Danzhking/secure-audit/services/processor/internal/repository"
	"github.com/Danzhking/secure-audit/services/processor/internal/service"
)

func main() {
	cfg := config.Load()

	db := repository.ConnectPostgres(cfg.PostgresURL)
	defer db.Close()

	repo := repository.NewEventRepository(db)
	if err := repo.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	conn := queue.ConnectRabbitMQ(cfg.RabbitURL)
	defer conn.Close()

	consumer, err := queue.NewConsumer(conn)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}
	defer consumer.Close()

	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatal("Failed to start consuming:", err)
	}

	eventService := service.NewEventService(repo)

	go eventService.ProcessMessages(msgs)

	log.Println("Processor started. Waiting for events...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Processor shutting down...")
}
