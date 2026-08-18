package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultDbUserValue     = "postgres"
	defaultDbPasswordValue = ""
	defaultDbHostValue     = "localhost"
	defaultDbPortValue     = "5432"
	defaultDbNameValue     = "postgres"

	defaultServerPort = ":8080"
)

type Config struct {
	DbOptions
	ServerOptions
}

type DbOptions struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

type ServerOptions struct {
	ListenAddr string
}

// Returns parsed connection string from config.
func (c *DbOptions) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
}

// Loads and returns config using environment variables.
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning .env file not found, using environment variables")
	}

	return &Config{
		DbOptions: DbOptions{
			DbUser:     getEnv("POSTGRES_USER", defaultDbUserValue),
			DbPassword: getEnv("POSTGRES_PASSWORD", defaultDbPasswordValue),
			DbHost:     getEnv("POSTGRES_HOST", defaultDbHostValue),
			DbPort:     getEnv("POSTGRES_PORT", defaultDbPortValue),
			DbName:     getEnv("POSTGRES_DB", defaultDbNameValue),
		},
		ServerOptions: ServerOptions{
			ListenAddr: getEnv("SERVER_PORT", defaultServerPort),
		},
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
