package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

type topology struct {
	RabbitMQ struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"rabbitmq"`
	Elasticsearch struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"elasticsearch"`
}

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

func loadTopology() topology {
	var t topology
	data, err := os.ReadFile("/shared/config/topology.json")
	if err != nil {
		log.Fatal("topology.json not found: ", err)
	}
	if err := json.Unmarshal(data, &t); err != nil {
		log.Fatal("failed to parse topology.json: ", err)
	}
	return t
}

func loadSharedConfig() sharedConfig {
	var cfg sharedConfig
	data, err := os.ReadFile("/shared/config/config.json")
	if err != nil {
		log.Fatal("config.json not found: ", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatal("failed to parse config.json: ", err)
	}
	return cfg
}

func LoadConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, reading from environment variables")
		}

		topo := loadTopology()
		shared := loadSharedConfig()

		instance = &Config{
			// Infrastructure — topology.json
			RabbitMQURL:       fmt.Sprintf("amqp://%s:%s@%s:%d/", topo.RabbitMQ.User, topo.RabbitMQ.Password, topo.RabbitMQ.Host, topo.RabbitMQ.Port),
			ElasticsearchURLs: []string{fmt.Sprintf("http://%s:%d", topo.Elasticsearch.Host, topo.Elasticsearch.Port)},

			// Secrets / service-specific — env var only
			OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
			WorkerCount: getEnvInt("WORKER_COUNT", 5),

			// Topology — config.json only
			RabbitMQQueue:             shared.Messaging.Queue,
			RawEventsExchange:         shared.Messaging.RawEventsExchange,
			ProcessedEventsExchange:   shared.Messaging.ProcessedEventsExchange,
			DLXExchange:               shared.Messaging.DLXExchange,
			DLXQueue:                  shared.Messaging.DLXQueue,
			ProcessedEventsIndex:      shared.Elasticsearch.Index,
			ElasticsearchGeocodeIndex: shared.Elasticsearch.GeocodeIndex,
			OpenAIModel:               shared.OpenAI.Model,
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
