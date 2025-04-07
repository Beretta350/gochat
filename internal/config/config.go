package config

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"

	"github.com/Beretta350/gochat/internal/app/adapters/wsadapter"
	"github.com/Beretta350/gochat/pkg/envutil"
	clientwrapper "github.com/Beretta350/gochat/pkg/kafkawrapper/client"
	consumerwrapper "github.com/Beretta350/gochat/pkg/kafkawrapper/consumer"
	producerwrapper "github.com/Beretta350/gochat/pkg/kafkawrapper/producer"
	"github.com/Beretta350/gochat/pkg/logger"
)

var (
	_, b, _, _  = runtime.Caller(0)
	internalDir = filepath.Dir(filepath.Dir(b))
	projectRoot = filepath.Dir(internalDir)
	configsDir  = filepath.Join(projectRoot, "configs")
)

func init() {
	env := envutil.GetEnv("ENV", "dev") // Default to "dev"
	if env == "local" {
		configPath := filepath.Join(configsDir, "local.env")
		if err := godotenv.Load(configPath); err != nil {
			panic(err)
		}
	}

	kafkaBrokers := envutil.GetEnv("KAFKA_BROKER", "kafka:9092")

	// Setup logger before all
	logger.Init(env)

	// Setup Kafka components
	clientwrapper.Init(kafkaBrokers)           // Admin client for topic management
	producerwrapper.InitProducer(kafkaBrokers) // Producer initialization
	consumerwrapper.InitConsumer(kafkaBrokers) // Consumer initialization

	setupApplication(env)

	logger.Info("Configuration successfully setup!")
}

var once sync.Once
var instance *Config

type Config struct {
	Server   *ServerConfig
	Upgrader wsadapter.Upgrader
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

func setupApplication(env string) *Config {
	once.Do(func() {
		serverConfig := &ServerConfig{
			Port: envutil.GetEnv("SERVER_PORT", "8080"),
		}

		var upgrader wsadapter.Upgrader
		if env != "prod" {
			upgrader = wsadapter.NewUpgrader(
				envutil.GetEnvInt("WS_READ_BUFFER_SIZE", 1024),
				envutil.GetEnvInt("WS_WRITE_BUFFER_SIZE", 1024),
				envutil.GetEnvBool("WS_CHECK_ORIGIN", true),
			)
		}

		instance = &Config{Server: serverConfig, Upgrader: upgrader}
	})
	return instance
}

func GetServerConfig() *ServerConfig {
	return instance.Server
}

func GetWebsocketUpgrader() wsadapter.Upgrader {
	return instance.Upgrader
}
