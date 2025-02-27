package config

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"

	"github.com/Beretta350/gochat/internal/config/wsupgrader"
	"github.com/Beretta350/gochat/pkg/kafka_wrapper"
	"github.com/Beretta350/gochat/pkg/logger"
	"github.com/Beretta350/gochat/pkg/util"
)

var (
	_, b, _, _  = runtime.Caller(0)
	internalDir = filepath.Dir(filepath.Dir(b))
	projectRoot = filepath.Dir(internalDir)
	configsDir  = filepath.Join(projectRoot, "configs")
)

func init() {
	env := util.GetEnv("ENV", "dev") // Default to "dev"
	if env == "local" {
		configPath := filepath.Join(configsDir, "local.env")
		if err := godotenv.Load(configPath); err != nil {
			panic(err)
		}
	}

	kafkaBrokers := util.GetEnv("KAFKA_BROKER", "kafka:9092")

	// Setup logger before all
	logger.Init(env)

	// Setup kafka admin client and wrapper
	kafka_wrapper.Init(kafkaBrokers)

	setupApplication(env)

	logger.Info("Configuration successfully setup!")
}

var once sync.Once
var instance *Config

type Config struct {
	Server   *ServerConfig
	Upgrader *websocket.Upgrader
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

func setupApplication(env string) *Config {
	once.Do(func() {
		serverConfig := &ServerConfig{
			Port: util.GetEnv("SERVER_PORT", "8080"),
		}

		var upgrader websocket.Upgrader
		if env != "prod" {
			upgrader = wsupgrader.NewUpgrader(
				util.GetEnvInt("WS_READ_BUFFER_SIZE", 1024),
				util.GetEnvInt("WS_WRITE_BUFFER_SIZE", 1024),
				util.GetEnvBool("WS_CHECK_ORIGIN", true),
			)
		}

		instance = &Config{Server: serverConfig, Upgrader: &upgrader}
	})
	return instance
}

func GetServerConfig() *ServerConfig {
	return instance.Server
}

func GetWebsocketUpgrader() *websocket.Upgrader {
	return instance.Upgrader
}
