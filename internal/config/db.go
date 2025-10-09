package config

import "github.com/SarkiMudboy/easebox-api/pkg/env"

type DBConfig struct {
	Addr            string
	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime int
}

func loadDBConfig() *DBConfig {
	return &DBConfig{
		Addr:            env.GetString("DB_ADDR", "admin:1234@/admin?parseTime=true"),
		MaxIdleConn:     env.GetInt("DB_MAX_IDLE_CONN", 10),
		MaxOpenConn:     env.GetInt("DB_MAX_OPEN_CONN", 10),
		MaxConnLifetime: env.GetInt("DB_MAX_CONN_LIFETIME", 10),
	}
}
