package main

import (
	"context"
	"log"
	"time"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	de "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	nominatim "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction/nominatimgeocoder"
	MessageQueue "github.com/IshaySela/israel-osint-ai/services/processing/messagebroker"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	"github.com/IshaySela/israel-osint-ai/services/processing/processor"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
	"github.com/IshaySela/israel-osint-ai/services/processing/workerpool"
	"golang.org/x/time/rate"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()
	rateLimiter := rate.NewLimiter(rate.Every(1100*time.Millisecond), 1)

	geocoder := de.NewGeocodingService(func(location string) (models.Geocode, *de.GeocodeError) {
		return nominatim.NominatimSearch(location, rateLimiter)
	})

	esClient := storage.NewElasticsearchClient()
	if err := esClient.Setup(cfg.ElasticsearchURLs); err != nil {
		log.Fatalf("Error setting up elasticsearch: %v", err)
	}

	pool := workerpool.NewWorkerPool(cfg.WorkerCount, 100)
	broker := MessageQueue.NewRabbitClient(cfg, pool)
	proc := processor.NewProcessor(cfg, geocoder, esClient, &broker)

	log.Printf("Starting %d workers...\n", cfg.WorkerCount)
	pool.Start(ctx, proc)

	log.Println("Starting message broker...")
	if err := broker.ListenForRawEvents(); err != nil {
		log.Fatalf("Error starting message broker: %v\n", err)
	}

	pool.Wait()
}
