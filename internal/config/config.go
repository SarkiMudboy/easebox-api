package config

import "github.com/SarkiMudboy/easebox-api/pkg/env"

type AppConfig struct {
	Port string
	ServerAddress string
}

func loadAppConfig() *AppConfig {
	return &AppConfig{
		Port: env.GetString("PORT", "8080"),
		ServerAddress: env.GetString("BASE_URL", "http://localhost"),
	}
}

type Config struct {
	App *AppConfig
	DB  *DBConfig
}

func Load() *Config {
	return &Config{
		App: loadAppConfig(),
		DB: loadDBConfig(),
	}
}