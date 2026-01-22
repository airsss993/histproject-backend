package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("ошибка подключения к PostgreSQL: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("ошибка пинга БД: %s", err.Error())
	}

	log.Println("Успешное подключение к БД!")

	return db
}
