package main

import (
	"context"
	"fmt"

	extractinfo "github.com/IshaySela/israel-osint-ai/services/processing/data-extraction"
	MessageQueue "github.com/IshaySela/israel-osint-ai/services/processing/messagebroker"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
)

func main() {
	broker := MessageQueue.NewRabbitListener("amqp://guest:guest@localhost:5672/", "osint_events")

	fmt.Println("Starting message broker...")
	done := make(chan bool)
	ctx := context.Background()
	geocoder := extractinfo.NewGeocodingService()

	err := broker.Listen(func(event models.RawOsintEvent) {
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

		fmt.Printf("Summary: %s\n", result.HeSummary)
		fmt.Println("Locations and Coordinates:")
		for i, location := range result.EnLocations {
			if i < len(coordinates) {
				fmt.Printf("- %s: Lat %s, Lon %s\n", location, coordinates[i].Lat, coordinates[i].Lon)
			} else {
				fmt.Printf("- %s: Coordinates not found\n", location)
			}
		}

	})

	if err != nil {
		fmt.Printf("Error starting message broker: %v\n", err)
	}

	<-done
}
