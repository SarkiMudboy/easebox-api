package config

import (
	"time"

	"github.com/SarkiMudboy/easebox-api/pkg/env"
)

type DBConfig struct {
	Addr            string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
}

func loadDBConfig() *DBConfig {
	return &DBConfig{
		Addr:            env.GetString("DB_ADDR", ""),
		MaxIdleConn:     env.GetInt("DB_MAX_IDLE_CONN", 10),
		MaxOpenConn:     env.GetInt("DB_MAX_OPEN_CONN", 10),
		MaxConnLifetime: time.Duration(env.GetInt("DB_MAX_CONN_LIFETIME", 10)) * time.Second,
	}
}
