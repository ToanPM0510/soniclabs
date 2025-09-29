package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr          string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration

	PGURL     string
	RedisAddr string
	RedisDB   int
}

func mustEnv(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		if def == "" {
			return def
		}
		log.Fatalf("missing env: %s", key)
	}
	return v
}
func mustEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("bad env %s: %v", key, err)
	}
	return i
}

func Load() Config {
	return Config{
		HTTPAddr:          mustEnv("HTTP_ADDR", ":8080"),
		ReadHeaderTimeout: 5 * time.Second,
		ShutdownTimeout:   10 * time.Second,
		PGURL:             mustEnv("PG_URL", ""),
		RedisAddr:         mustEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:           mustEnvInt("REDIS_DB", 0),
	}
}
