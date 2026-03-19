package main

import (
	"context"
	"log"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	dataextraction "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction"
	nominatim "github.com/IshaySela/israel-osint-ai/services/processing/dataextraction/nominatimgeocoder"
	MessageQueue "github.com/IshaySela/israel-osint-ai/services/processing/messagebroker"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
)

func main() {
	cfg := config.LoadConfig()

	broker := MessageQueue.NewRabbitListener(cfg.RabbitMQURL, cfg.RabbitMQQueue)

	log.Println("Starting message broker...")
	done := make(chan bool)
	ctx := context.Background()
	geocoder := dataextraction.NewGeocodingService(nominatim.NominatimSearch)

	esClient := storage.NewElasticsearchClient()
	err := esClient.Setup(cfg.ElasticsearchURLs)
	if err != nil {
		log.Fatalf("Error setting up elasticsearch: %v", err)
	}

	err = broker.Listen(func(event models.RawOsintEvent) {
		log.Printf("Received event: %s\n", string(event.Text))
		result, err := dataextraction.CreateAgentSummary(event, ctx, cfg.OpenAIKey, cfg.OpenAIModel)

		if err != nil {
			log.Printf("Error extracting info: %v\n", err)
			return
		}

		log.Printf("AI Summary: %+v\n", result)

		coordinates, geocodeErr := geocoder.GetBatchCoordinates(result.EnLocations)
		if geocodeErr != nil {
			log.Printf("Error fetching coordinates: %v\n", err)
			return
		}

		locationMap := make(map[string]dataextraction.Geocode)

		for i, location := range result.EnLocations {
			if i < len(coordinates) {
				locationMap[location] = coordinates[i]
			} else {
				log.Printf("- %s: Coordinates not found\n", location)
			}
		}

		processedEvent := storage.ProcessedEvent{
			RawMessage: event.Text,
			Summary:    result.HeSummary,
			Locations:  locationMap,
			Timestamp:  event.Date,
		}

		err = esClient.IndexEvent(ctx, cfg.ElasticsearchIndex, processedEvent)
		if err != nil {
			log.Printf("Error indexing event to elasticsearch: %v\n", err)
		} else {
			log.Println("Successfully indexed event to elasticsearch")
		}

	})

	if err != nil {
		log.Printf("Error starting message broker: %v\n", err)
	}

	<-done
}
