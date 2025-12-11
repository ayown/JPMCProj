package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
	ML       MLConfig
	Server   ServerConfig
}

type AppConfig struct {
	Env      string
	Name     string
	LogLevel string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConnections  int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	CacheTTL time.Duration
}

type KafkaConfig struct {
	Brokers              []string
	TopicVerification    string
	TopicReports         string
	TopicAlerts          string
	ConsumerGroup        string
	AutoOffsetReset      string
	SessionTimeout       time.Duration
	HeartbeatInterval    time.Duration
}

type JWTConfig struct {
	Secret               string
	Expiry               time.Duration
	RefreshTokenExpiry   time.Duration
}

type MLConfig struct {
	ServiceURL       string
	InferenceTimeout time.Duration
}

type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (for local development)
	_ = godotenv.Load()

	config := &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "development"),
			Name:     getEnv("APP_NAME", "fraud-detection-system"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DATABASE_HOST", "localhost"),
			Port:            getEnvAsInt("DATABASE_PORT", 5432),
			User:            getEnv("DATABASE_USER", "frauddetection"),
			Password:        getEnv("DATABASE_PASSWORD", "frauddetection_password"),
			DBName:          getEnv("DATABASE_NAME", "frauddetection_db"),
			SSLMode:         getEnv("DATABASE_SSL_MODE", "disable"),
			MaxConnections:  getEnvAsInt("DATABASE_MAX_CONNECTIONS", 100),
			MaxIdleConns:    getEnvAsInt("DATABASE_MAX_IDLE_CONNECTIONS", 10),
			ConnMaxLifetime: time.Hour,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			CacheTTL: getEnvAsDuration("REDIS_CACHE_TTL", 3600*time.Second),
		},
		Kafka: KafkaConfig{
			Brokers:              []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			TopicVerification:    getEnv("KAFKA_TOPIC_VERIFICATION", "verification-requests"),
			TopicReports:         getEnv("KAFKA_TOPIC_REPORTS", "fraud-reports"),
			TopicAlerts:          getEnv("KAFKA_TOPIC_ALERTS", "fraud-alerts"),
			ConsumerGroup:        getEnv("KAFKA_CONSUMER_GROUP", "fraud-detection-workers"),
			AutoOffsetReset:      getEnv("KAFKA_AUTO_OFFSET_RESET", "earliest"),
			SessionTimeout:       30 * time.Second,
			HeartbeatInterval:    3 * time.Second,
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "your-secret-key"),
			Expiry:             getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
			RefreshTokenExpiry: getEnvAsDuration("REFRESH_TOKEN_EXPIRY", 168*time.Hour),
		},
		ML: MLConfig{
			ServiceURL:       getEnv("ML_SERVICE_URL", "http://localhost:8000"),
			InferenceTimeout: getEnvAsDuration("INFERENCE_TIMEOUT", 5*time.Second),
		},
		Server: ServerConfig{
			Port:            getEnv("API_GATEWAY_PORT", "8080"),
			Host:            getEnv("API_GATEWAY_HOST", "0.0.0.0"),
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 30 * time.Second,
		},
	}

	// Validate required fields
	if config.JWT.Secret == "your-secret-key" && config.App.Env == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return config, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

