package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL               string
	RabbitMQQueue             string
	RawEventsExchange         string
	ProcessedEventsExchange   string
	DLXExchange               string
	DLXQueue                  string
	ElasticsearchURLs         []string
	ProcessedEventsIndex      string
	ElasticsearchGeocodeIndex string
	OpenAIKey                 string
	OpenAIModel               string
	WorkerCount               int
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, reading from environment variables")
		}

		instance = &Config{
			RabbitMQURL:               getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			RabbitMQQueue:             getEnv("RABBITMQ_QUEUE", "osint_events"),
			RawEventsExchange:         getEnv("RAW_EVENTS_EXCHANGE", "raw_events"),
			ProcessedEventsExchange:   getEnv("PROCESSED_EVENTS_EXCHANGE", "processed_events"),
			DLXExchange:               getEnv("DLX_EXCHANGE", "dead_letter"),
			DLXQueue:                  getEnv("DLX_QUEUE", "dead_letter_queue"),
			ElasticsearchURLs:         strings.Split(getEnv("ELASTICSEARCH_URLS", "http://localhost:9200"), ","),
			ProcessedEventsIndex:      getEnv("ELASTICSEARCH_INDEX", "osint_events"),
			ElasticsearchGeocodeIndex: getEnv("ELASTICSEARCH_GEOCODE_INDEX", "geocode_cache"),
			OpenAIKey:                 getEnv("OPENAI_API_KEY", ""),
			OpenAIModel:               getEnv("OPENAI_MODEL", "gpt-5-mini"),
			WorkerCount:               getEnvInt("WORKER_COUNT", 5),
		}
	})

	return instance
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var result int
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
	}
	return defaultValue
}
