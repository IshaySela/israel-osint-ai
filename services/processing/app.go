package main

import (
	"context"
	"fmt"
	"log"

	extractinfo "github.com/IshaySela/israel-osint-ai/services/processing/data-extraction"
	MessageQueue "github.com/IshaySela/israel-osint-ai/services/processing/messagebroker"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
)

func main() {
	broker := MessageQueue.NewRabbitListener("amqp://guest:guest@localhost:5672/", "osint_events")

	fmt.Println("Starting message broker...")
	done := make(chan bool)
	ctx := context.Background()
	geocoder := extractinfo.NewGeocodingService()

	esClient := storage.NewElasticsearchClient()
	err := esClient.Setup([]string{"http://localhost:9200"})
	if err != nil {
		log.Fatalf("Error setting up elasticsearch: %v", err)
	}

	err = broker.Listen(func(event models.RawOsintEvent) {
		fmt.Printf("Received event: %s\n", string(event.Text))
		result, err := extractinfo.CreateAgentSummary(event, ctx)

		if err != nil {
			fmt.Printf("Error extracting info: %v\n", err)
			return
		}
		coordinates, err := geocoder.GetBatchCoordinates(result.EnLocations)
		if err != nil {
			fmt.Printf("Error fetching coordinates: %v\n", err)
			return
		}

		locationMap := make(map[string]extractinfo.Geocode)
		fmt.Printf("Summary: %s\n", result.HeSummary)
		fmt.Println("Locations and Coordinates:")
		for i, location := range result.EnLocations {
			if i < len(coordinates) {
				fmt.Printf("- %s: Lat %s, Lon %s\n", location, coordinates[i].Lat, coordinates[i].Lon)
				locationMap[location] = coordinates[i]
			} else {
				fmt.Printf("- %s: Coordinates not found\n", location)
			}
		}

		processedEvent := storage.ProcessedEvent{
			RawMessage: event.Text,
			Summary:    result.HeSummary,
			Locations:  locationMap,
			Timestamp:  event.Date,
		}

		err = esClient.IndexEvent(ctx, "osint_events", processedEvent)
		if err != nil {
			fmt.Printf("Error indexing event to elasticsearch: %v\n", err)
		} else {
			fmt.Println("Successfully indexed event to elasticsearch")
		}

	})

	if err != nil {
		fmt.Printf("Error starting message broker: %v\n", err)
	}

	<-done
}
