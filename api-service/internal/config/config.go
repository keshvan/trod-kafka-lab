package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerHost string
	ServerPort string

	DataClientHost string
	DataClientPort string
	KafkaHost      string
	KafkaPort      string
	KafkaTopic     string
}

func MustLoad() *Config {
	return &Config{
		ServerHost:     getEnv("SERVER_HOST", "localhost"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DataClientHost: requireEnv("DATA_CLIENT_HOST"),
		DataClientPort: requireEnv("DATA_CLIENT_PORT"),
		KafkaHost:      getEnv("KAFKA_HOST", ""),
		KafkaPort:      getEnv("KAFKA_PORT", ""),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "appointments"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}
