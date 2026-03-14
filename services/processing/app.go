package main

import (
	"fmt"

	MessageQueue "github.com/IshaySela/israel-osint-ai/services/processing/messagebroker"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
)

func main() {
	broker := MessageQueue.NewRabbitListener("amqp://guest:guest@localhost:5672/", "osint_events")

	fmt.Println("Starting message broker...")
	done := make(chan bool)

	err := broker.Listen(func(event models.RawOsintEvent) {
		fmt.Printf("Received event: %s\n", string(event.Text))
	})

	if err != nil {
		fmt.Printf("Error starting message broker: %v\n", err)
	}

	<-done
}
