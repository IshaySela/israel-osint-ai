package config

import (
	"encoding/json"
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

type sharedConfig struct {
	Messaging struct {
		Queue                   string `json:"queue"`
		RawEventsExchange       string `json:"raw_events_exchange"`
		ProcessedEventsExchange string `json:"processed_events_exchange"`
		DLXExchange             string `json:"dlx_exchange"`
		DLXQueue                string `json:"dlx_queue"`
	} `json:"messaging"`
	Elasticsearch struct {
		Index        string `json:"index"`
		GeocodeIndex string `json:"geocode_index"`
	} `json:"elasticsearch"`
	OpenAI struct {
		Model string `json:"model"`
	} `json:"openai"`
}

func loadSharedConfig() sharedConfig {
	var cfg sharedConfig
	data, err := os.ReadFile("/shared/config/config.json")
	if err != nil {
		log.Println("No shared config.json found, using defaults")
		return cfg
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("Failed to parse shared config.json: %v", err)
	}
	return cfg
}

func LoadConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, reading from environment variables")
		}

		shared := loadSharedConfig()

		instance = &Config{
			RabbitMQURL:               getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			RabbitMQQueue:             getEnv("RABBITMQ_QUEUE", shared.Messaging.Queue),
			RawEventsExchange:         getEnv("RAW_EVENTS_EXCHANGE", shared.Messaging.RawEventsExchange),
			ProcessedEventsExchange:   getEnv("PROCESSED_EVENTS_EXCHANGE", shared.Messaging.ProcessedEventsExchange),
			DLXExchange:               getEnv("DLX_EXCHANGE", shared.Messaging.DLXExchange),
			DLXQueue:                  getEnv("DLX_QUEUE", shared.Messaging.DLXQueue),
			ElasticsearchURLs:         strings.Split(getEnv("ELASTICSEARCH_URLS", "http://localhost:9200"), ","),
			ProcessedEventsIndex:      getEnv("ELASTICSEARCH_INDEX", shared.Elasticsearch.Index),
			ElasticsearchGeocodeIndex: getEnv("ELASTICSEARCH_GEOCODE_INDEX", shared.Elasticsearch.GeocodeIndex),
			OpenAIKey:                 getEnv("OPENAI_API_KEY", ""),
			OpenAIModel:               getEnv("OPENAI_MODEL", shared.OpenAI.Model),
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
