package config

import (
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/util"
	"github.com/joho/godotenv"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

func init() {
	env := util.GetEnv("ENV", "dev") // Default to "dev"
	if env == "local" {
		configPath := filepath.Join(basePath, "local.env")
		if err := godotenv.Load(configPath); err != nil {
			panic(err)
		}
	}

	logger.Init(env)
	setup()

	logger.Info("Configuration successfully setup!")
}

var once sync.Once
var instance *Config

type Config struct {
	Server *ServerConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

func setup() *Config {
	once.Do(func() {
		serverConfig := &ServerConfig{
			Port: util.GetEnv("SERVER_PORT", "8080"),
		}
		instance = &Config{Server: serverConfig}
	})
	return instance
}

func GetServerConfig() *ServerConfig {
	return instance.Server
}
