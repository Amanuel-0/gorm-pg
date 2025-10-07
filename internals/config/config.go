package config

import (
	"os"
)

type (
	Container struct {
		AppConfig *App
		DB        *DB
	}

	App struct {
		ServerDomain string
	}

	DB struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
	}
)

func New() (*Container, error) {
	// Initialize the application configuration
	app := &App{}

	// Initialize the database configuration
	db := &DB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	return &Container{app, db}, nil
}

// getEnvValue returns the environment variable value for key, or dv if unset or empty.
func getEnvValue(key, dv string) string {
	value := os.Getenv(key)
	if value == "" {
		return dv
	}
	return value
}
