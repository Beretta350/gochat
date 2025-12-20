package config

import (
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"

	"github.com/Beretta350/gochat/pkg/envutil"
	"github.com/Beretta350/gochat/pkg/logger"
)

var (
	_, b, _, _  = runtime.Caller(0)
	internalDir = filepath.Dir(filepath.Dir(b))
	projectRoot = filepath.Dir(internalDir)
	configsDir  = filepath.Join(projectRoot, "configs")
)

// Config holds all application configuration
type Config struct {
	Env    string
	Server ServerConfig
	Redis  RedisConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// NewConfig creates a new Config (Fx provider)
func NewConfig() (*Config, error) {
	env := envutil.GetEnv("ENV", "dev")

	// Load .env file for local environment
	if env == "local" {
		configPath := filepath.Join(configsDir, "local.env")
		if err := godotenv.Load(configPath); err != nil {
			return nil, err
		}
	}

	// Initialize logger
	logger.Init(env)

	cfg := &Config{
		Env: env,
		Server: ServerConfig{
			Port: envutil.GetEnv("SERVER_PORT", "8080"),
		},
		Redis: RedisConfig{
			Addr:     envutil.GetEnv("REDIS_ADDR", "localhost:6379"),
			Password: envutil.GetEnv("REDIS_PASSWORD", ""),
			DB:       envutil.GetEnvInt("REDIS_DB", 0),
		},
	}

	logger.Info("Configuration loaded")
	return cfg, nil
}
