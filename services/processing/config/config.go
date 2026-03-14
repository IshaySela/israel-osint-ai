package config

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL        string
	RabbitMQQueue      string
	ElasticsearchURLs  []string
	ElasticsearchIndex string
	OpenAIKey          string
	OpenAIModel        string
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
			RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			RabbitMQQueue:      getEnv("RABBITMQ_QUEUE", "osint_events"),
			ElasticsearchURLs:  strings.Split(getEnv("ELASTICSEARCH_URLS", "http://localhost:9200"), ","),
			ElasticsearchIndex: getEnv("ELASTICSEARCH_INDEX", "osint_events"),
			OpenAIKey:          getEnv("OPENAI_API_KEY", ""),
			OpenAIModel:        getEnv("OPENAI_MODEL", "gpt-5-mini"),
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
