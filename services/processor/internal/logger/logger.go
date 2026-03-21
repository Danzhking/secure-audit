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
		log.Fatal("Failed to initialize logger:", err)
	}

	zapLogger = zapLogger.With(zap.String("service", serviceName))
	zap.ReplaceGlobals(zapLogger)
}
