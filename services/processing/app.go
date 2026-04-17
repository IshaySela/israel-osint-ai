package main

import (
	"context"
	"log"
	"time"

	"processing/config"
	de "processing/dataextraction"
	MessageQueue "processing/messagebroker"
	"processing/processor"
	storage "processing/storage"
	"processing/workerpool"
	"golang.org/x/time/rate"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()
	rateLimiter := rate.NewLimiter(rate.Every(1100*time.Millisecond), 1)

	esClient := storage.NewElasticsearchClient(cfg)
	if err := esClient.Setup(cfg.ElasticsearchURLs); err != nil {
		log.Fatalf("Error setting up elasticsearch: %v", err)
	}

	geocoder, err := de.NewGeocodingServiceBuilder().
		WithContext(ctx).
		WithNominatim(rateLimiter).
		WithElasticsearchCache(esClient).
		Build()
	if err != nil {
		log.Fatalf("Error building geocoding service: %v", err)
	}

	pool := workerpool.NewWorkerPool(cfg.WorkerCount, 100)
	broker := MessageQueue.NewRabbitClient(cfg, pool)
	proc := processor.NewProcessor(cfg, geocoder, esClient, &broker)

	log.Printf("Starting %d workers...\n", cfg.WorkerCount)
	pool.Start()

	log.Println("Starting message broker...")
	if err := broker.ListenForRawEvents(ctx, proc); err != nil {
		log.Fatalf("Error starting message broker: %v\n", err)
	}

	pool.Wait()
}
