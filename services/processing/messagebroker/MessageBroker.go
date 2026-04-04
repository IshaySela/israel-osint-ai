package messagebroker

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	"github.com/IshaySela/israel-osint-ai/services/processing/storage"
	"github.com/IshaySela/israel-osint-ai/services/processing/workerpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventProcessor interface {
	Process(ctx context.Context, event models.RawOsintEvent) error
}

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
	config  *config.Config
	pool    *workerpool.WorkerPool
}

func NewRabbitClient(config *config.Config, pool *workerpool.WorkerPool) RabbitClient {
	return RabbitClient{config: config, pool: pool}
}

func (rl *RabbitClient) setup() error {
	conn, err := amqp.Dial(rl.config.RabbitMQURL)
	if err != nil {
		return errors.New("failed to establish connection to rabbitmq host")
	}
	rl.conn = conn
	ch, err := rl.conn.Channel()

	if err != nil {
		return errors.New("failed to open channel to rabbitmq host")
	}

	err = ch.ExchangeDeclare(rl.config.ProcessedEventsExchange,
		"fanout",
		true,  // durable
		false, false, false, nil)

	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		rl.config.RabbitMQQueue, // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)

	if err != nil {
		return errors.New("failed to declare queue")
	}

	rl.channel = ch
	rl.queue = &q
	return nil
}

func (rl *RabbitClient) ListenForRawEvents(ctx context.Context, proc EventProcessor) error {
	if err := rl.setup(); err != nil {
		return err
	}

	msgs, err := rl.channel.Consume(
		rl.queue.Name, // queue
		"",            // consumer
		false,         // auto-ack = false
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)

	if err != nil {
		return errors.New("failed to register a consumer")
	}

	go func() {
		for d := range msgs {
			d := d
			var event models.RawOsintEvent
			if err := event.Unmarshal(d.Body); err != nil {
				log.Printf("Failed to unmarshal message, discarding: %v", err)
				d.Nack(false, false)
				continue
			}
			rl.pool.Submit(func() {
				if err := proc.Process(ctx, event); err != nil {
					log.Printf("Failed to process event, requeueing: %v", err)
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}
			})
		}
	}()

	return nil
}

func (rl *RabbitClient) Publish(exchange string, routingKey string, body []byte) error {
	if rl.channel == nil {
		return errors.New("failed to publish message, channel not initialized")
	}

	return rl.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (rl *RabbitClient) PublishProcessedEvent(ev storage.ProcessedEvent, dbId string) error {
	msg := models.ProcessedEventMessage{
		DbId:      dbId,
		Summary:   ev.Summary,
		Locations: ev.Locations,
		Timestamp: ev.Timestamp,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return rl.Publish(rl.config.ProcessedEventsExchange, "", body)
}
