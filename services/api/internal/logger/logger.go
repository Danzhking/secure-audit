package logger

import (
	"log"

	"go.uber.org/zap"
)

func Init(serviceName string) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}

	zapLogger, err := cfg.Build()
	if err != nil {
		log.Fatal("Не удалось инициализировать логгер:", err)
	}

	zapLogger = zapLogger.With(zap.String("service", serviceName))
	zap.ReplaceGlobals(zapLogger)
}
