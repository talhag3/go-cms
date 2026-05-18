package config

import (
	"fmt"
	"os"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Auth     AuthConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name  string
	Env   string
	Debug bool
}

type ServerConfig struct {
	Host string
	Port string
}

type AuthConfig struct {
	SecretKey   string
	TokenExpiry int // in hours
}

type DatabaseConfig struct {
	Host     string // Database host (default: localhost)
	Port     string // Database port (default: 5432)
	User     string // Database user
	Password string // Database password
	DBName   string // Database name
	SSLMode  string // SSL mode: disable, require, verify-ca, verify-full
	// Pool settings
	MaxConns    int32 // Maximum connections in pool
	MinConns    int32 // Minimum connections in pool (kept alive)
	MaxIdleTime int   // Max time a connection can sit idle (seconds)
	MaxLifetime int   // Max time a connection can live (seconds)
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:  getEnv("APP_NAME", "CMS API"),
			Env:   getEnv("APP_ENV", "development"),
			Debug: getEnv("APP_DEBUG", "true") == "true",
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "3000"),
		},
		Auth: AuthConfig{
			SecretKey:   getEnv("AUTH_SECRET", "qwert1234"),
			TokenExpiry: 24,
		},
		Database: DatabaseConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5432"),
			User:        getEnv("DB_USER", "capt41"),
			Password:    getEnv("DB_PASSWORD", "pass1234"),
			DBName:      getEnv("DB_NAME", "gocms"),
			SSLMode:     getEnv("DB_SSL_MODE", "disable"),
			MaxConns:    25,   // Good default for most apps
			MinConns:    5,    // Keep 5 connections alive
			MaxIdleTime: 300,  // 5 minutes
			MaxLifetime: 1800, // 30 minutes
		},
	}
}

// DSN returns the PostgreSQL connection string
//
// GO CONCEPT - STRING FORMATTING:
// fmt.Sprintf is like PHP's sprintf()
// DSN = Data Source Name (standard connection string format)
//
// PostgreSQL DSN format:
// postgres://user:password@host:port/dbname?sslmode=disable
//
// In PHP (Laravel):
// DB_CONNECTION=pgsql
// DB_HOST=127.0.0.1
// DB_PORT=5432
// DB_DATABASE=laravel
// DB_USERNAME=root
// DB_PASSWORD=secret
// (Laravel builds DSN internally)
func (dc *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DBName,
		dc.SSLMode,
	)
}

// Note: lowercase "getEnv" means this is PRIVATE (unexported)
// Can only be used within this package
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
