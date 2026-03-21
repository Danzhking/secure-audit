package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func ConnectPostgres(url string) *sql.DB {
	var db *sql.DB
	var err error

	for i := range 10 {
		db, err = sql.Open("postgres", url)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				zap.L().Info("PostgreSQL connected")
				return db
			}
		}
		zap.L().Warn("PostgreSQL not ready",
			zap.Int("attempt", i+1),
			zap.Error(err),
		)
		time.Sleep(2 * time.Second)
	}

	zap.L().Fatal("PostgreSQL connection failed after 10 attempts", zap.Error(err))
	return nil
}
