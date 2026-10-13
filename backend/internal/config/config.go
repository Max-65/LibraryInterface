package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	AllowedOrigin string
}

func Load() Config {
	return Config{
		Port:          getenv("PORT", "8080"),
		DatabaseURL:   getenv("DATABASE_URL", ""),
		AllowedOrigin: getenv("ALLOWED_ORIGIN", "http://localhost:5173"),
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
