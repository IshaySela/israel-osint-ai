package main

import (
	"context"
	"log"
	"sync"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	dataextraction "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	nominatim "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction/nominatimgeocoder"
	MessageQueue "github.com/IshaySela/israel-osint-ai/services/processing/messagebroker"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	"github.com/IshaySela/israel-osint-ai/services/processing/processor"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
)

func main() {
	cfg := config.LoadConfig()
	wg := sync.WaitGroup{}

	broker := MessageQueue.NewRabbitListener(cfg.RabbitMQURL, cfg.RabbitMQQueue)

	log.Println("Starting message broker...")
	ctx := context.Background()
	geocoder := dataextraction.NewGeocodingService(nominatim.NominatimSearch)

	esClient := storage.NewElasticsearchClient()
	err := esClient.Setup(cfg.ElasticsearchURLs)
	if err != nil {
		log.Fatalf("Error setting up elasticsearch: %v", err)
	}

	proc := processor.NewProcessor(cfg, geocoder, esClient)
	taskQueue := make(chan models.RawOsintEvent, 100)

	log.Printf("Starting %d workers...\n", cfg.WorkerCount)
	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)
		go func() {
			proc.StartWorker(ctx, taskQueue)
			wg.Done()
		}()
	}

	err = broker.Listen(func(event models.RawOsintEvent) {
		taskQueue <- event
	})

	if err != nil {
		log.Printf("Error starting message broker: %v\n", err)
	}

	wg.Wait()
}
