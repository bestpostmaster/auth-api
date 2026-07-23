package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Address       string
	MySQLDSN      string
	AllowedOrigin string
	JWTPrivateKey []byte
	JWTTTL        time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Address:       envOrDefault("HTTP_ADDRESS", ":8080"),
		MySQLDSN:      os.Getenv("MYSQL_DSN"),
		AllowedOrigin: envOrDefault("ALLOWED_ORIGIN", "http://localhost:4201"),
	}
	if cfg.MySQLDSN == "" {
		return Config{}, errors.New("MYSQL_DSN is required")
	}

	ttlMinutes, err := strconv.Atoi(envOrDefault("JWT_TTL_MINUTES", "60"))
	if err != nil || ttlMinutes <= 0 {
		return Config{}, errors.New("JWT_TTL_MINUTES must be a positive integer")
	}
	cfg.JWTTTL = time.Duration(ttlMinutes) * time.Minute

	keyText := strings.TrimSpace(os.Getenv("JWT_PRIVATE_KEY"))
	if keyText != "" {
		cfg.JWTPrivateKey = []byte(strings.ReplaceAll(keyText, `\n`, "\n"))
		return cfg, nil
	}
	keyFile := os.Getenv("JWT_PRIVATE_KEY_FILE")
	if keyFile == "" {
		return Config{}, errors.New("JWT_PRIVATE_KEY or JWT_PRIVATE_KEY_FILE is required")
	}
	cfg.JWTPrivateKey, err = os.ReadFile(keyFile)
	if err != nil {
		return Config{}, fmt.Errorf("read JWT_PRIVATE_KEY_FILE: %w", err)
	}
	return cfg, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
