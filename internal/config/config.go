package config

import "os"

type Config struct {
	App    AppConfig
	Server ServerConfig
	Auth   AuthConfig
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
	}
}

// Note: lowercase "getEnv" means this is PRIVATE (unexported)
// Can only be used within this package
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
