package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/SarkiMudboy/easebox-api/internal/config"

	_ "github.com/lib/pq"
)


func New(cfg *config.DBConfig) (*sql.DB, error) {
	
	var db *sql.DB
	var err error
	
	db, err = sql.Open("postgres", cfg.Addr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
