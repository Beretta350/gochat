package config

import (
	"path/filepath"
	"runtime"
	"time"

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
	Env      string
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	CORS     CORSConfig
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	URL string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	TLS      bool // Enable TLS for cloud Redis (Upstash, etc.)
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Domain   string
	Secure   bool   // true for HTTPS
	SameSite string // "Strict", "Lax", or "None"
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
		Database: DatabaseConfig{
			URL: envutil.GetEnv("DATABASE_URL", "postgres://gochat:gochat@localhost:5432/gochat?sslmode=disable"),
		},
		Redis: RedisConfig{
			Addr:     envutil.GetEnv("REDIS_ADDR", "localhost:6379"),
			Password: envutil.GetEnv("REDIS_PASSWORD", ""),
			DB:       envutil.GetEnvInt("REDIS_DB", 0),
			TLS:      envutil.GetEnvBool("REDIS_TLS", false), // true for Upstash/cloud Redis
		},
		JWT: JWTConfig{
			Secret:        envutil.GetEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
			AccessExpiry:  envutil.GetEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: envutil.GetEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		Cookie: CookieConfig{
			Domain:   envutil.GetEnv("COOKIE_DOMAIN", "localhost"),
			Secure:   envutil.GetEnvBool("COOKIE_SECURE", false), // true em produção (HTTPS)
			SameSite: envutil.GetEnv("COOKIE_SAMESITE", "Lax"),
		},
		CORS: CORSConfig{
			AllowedOrigins: envutil.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		},
	}

	logger.Info("Configuration loaded")
	return cfg, nil
}
