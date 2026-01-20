package db

import (
	"database/sql"
	"log"

	"github.com/airsss993/histproject-backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %s", err.Error())
	}

	log.Println("Successfully connected to DB!")

	return db
}
