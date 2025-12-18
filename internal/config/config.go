package config

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"

	"github.com/Beretta350/gochat/pkg/envutil"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/redisclient"
)

var (
	_, b, _, _  = runtime.Caller(0)
	internalDir = filepath.Dir(filepath.Dir(b))
	projectRoot = filepath.Dir(internalDir)
	configsDir  = filepath.Join(projectRoot, "configs")
)

func init() {
	env := envutil.GetEnv("ENV", "dev")
	if env == "local" {
		configPath := filepath.Join(configsDir, "local.env")
		if err := godotenv.Load(configPath); err != nil {
			panic(err)
		}
	}

	// Setup logger
	logger.Init(env)

	// Setup Redis
	redisConfig := redisclient.Config{
		Addr:     envutil.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: envutil.GetEnv("REDIS_PASSWORD", ""),
		DB:       envutil.GetEnvInt("REDIS_DB", 0),
	}
	if err := redisclient.Init(redisConfig); err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Setup application
	setupApplication(env)

	logger.Info("Configuration successfully setup!")
}

var appOnce sync.Once
var appInstance *AppConfig

type AppConfig struct {
	Server *ServerConfig
	Env    string
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

func setupApplication(env string) *AppConfig {
	appOnce.Do(func() {
		serverConfig := &ServerConfig{
			Port: envutil.GetEnv("SERVER_PORT", "8080"),
		}

		appInstance = &AppConfig{
			Server: serverConfig,
			Env:    env,
		}
	})
	return appInstance
}

func GetServerConfig() *ServerConfig {
	return appInstance.Server
}

func GetEnv() string {
	return appInstance.Env
}
