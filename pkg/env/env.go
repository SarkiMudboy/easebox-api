package env

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	env := dir(".env.local")
	err := godotenv.Load(env)
	if err != nil {
		log.Fatalf("Error loading environment variables: %s", err.Error())
	}
}

// from "https://github.com/joho/godotenv/issues/126"
// Finds the env file which typically stays in the root where the
// go.mod file is located, and returns its abs path
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			log.Panic("go.mod not found")
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}

func GetString(variable string, fallback string) string {
	value, ok := os.LookupEnv(variable)
	if !ok {
		return fallback
	}

	return value
}

func GetInt(variable string, fallback int) int {
	value, ok := os.LookupEnv(variable)
	if !ok {
		return fallback
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intVal
}
