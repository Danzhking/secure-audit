package repository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func ConnectPostgres(url string) *sql.DB {
	var db *sql.DB
	var err error

	for i := range 10 {
		db, err = sql.Open("postgres", url)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				log.Println("PostgreSQL connected")
				return db
			}
		}
		log.Printf("PostgreSQL not ready (attempt %d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("PostgreSQL connection failed after 10 attempts:", err)
	return nil
}
