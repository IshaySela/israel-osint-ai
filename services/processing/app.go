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

	err := broker.Listen(func(event models.RawOsintEvent) {
		fmt.Printf("Received event: %s\n", string(event.Text))
		result, err := extractinfo.CreateAgentSummary(event, ctx)

		if err != nil {
			fmt.Printf("Error extracting info: %v\n", err)
		} else {
			fmt.Printf("Sumamry: %s\n", result)
		}
	})

	if err != nil {
		fmt.Printf("Error starting message broker: %v\n", err)
	}

	<-done
}
