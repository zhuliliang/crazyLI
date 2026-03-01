package config

import (
    "log"
    "os"
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
    Port               string
    DatabaseURL        string
    GoogleClientID     string
    GoogleClientSecret string
    GoogleRedirectURL  string
    FrontendURL        string
    JWTSecret          string
}

// Load reads configuration from environment variables with sane defaults.
func Load() Config {
    cfg := Config{
        Port:               getEnv("PORT", "8080"),
        DatabaseURL:        os.Getenv("DATABASE_URL"),
        GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
        FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
        JWTSecret:          os.Getenv("JWT_SECRET"),
    }

    missing := []string{}
    if cfg.DatabaseURL == "" {
        missing = append(missing, "DATABASE_URL")
    }
    if cfg.GoogleClientID == "" {
        missing = append(missing, "GOOGLE_CLIENT_ID")
    }
    if cfg.GoogleClientSecret == "" {
        missing = append(missing, "GOOGLE_CLIENT_SECRET")
    }
    if cfg.GoogleRedirectURL == "" {
        missing = append(missing, "GOOGLE_REDIRECT_URL")
    }
    if cfg.JWTSecret == "" {
        missing = append(missing, "JWT_SECRET")
    }

    if len(missing) > 0 {
        log.Printf("[config] warning: missing env vars: %v", missing)
    }

    return cfg
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}
