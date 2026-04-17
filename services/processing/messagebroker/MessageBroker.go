package messagebroker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"processing/config"
	models "processing/models"
	"processing/storage"
	"processing/workerpool"
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
			event, err := models.ParseRawOsintEvent(d.Body)
			if err != nil {
				log.Printf("Failed to unmarshal message, discarding: %v", err)
				d.Nack(false, false)
				continue
			}
			rl.pool.Submit(func() {
				if err := proc.Process(ctx, event); err != nil {
					log.Printf("Failed to process event: %v", err)
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

func (rl *RabbitClient) PublishProcessedEvent(ev storage.ProcessedEvent[any], dbId string) error {
	msg := CreateMessageFromEvent(ev, dbId)

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return rl.Publish(rl.config.ProcessedEventsExchange, "", body)
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

	// Declare DLX exchange and queue
	if err = ch.ExchangeDeclare(rl.config.DLXExchange, "fanout", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare DLX exchange: %w", err)
	}
	log.Printf("Declared DLX exchange: %s", rl.config.DLXExchange)

	dlxQueue, err := ch.QueueDeclare(rl.config.DLXQueue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare DLX queue: %w", err)
	}
	log.Printf("Declared DLX queue: %s", dlxQueue.Name)

	if err = ch.QueueBind(dlxQueue.Name, "#", rl.config.DLXExchange, false, nil); err != nil {
		return fmt.Errorf("failed to bind DLX queue: %w", err)
	}
	log.Printf("Bound DLX queue %s to exchange %s", dlxQueue.Name, rl.config.DLXExchange)

	// Declare raw events fanout exchange
	if err = ch.ExchangeDeclare(rl.config.RawEventsExchange, "fanout", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare raw events exchange: %w", err)
	}
	log.Printf("Declared raw events exchange: %s", rl.config.RawEventsExchange)

	// Declare raw events queue with DLX routing
	q, err := ch.QueueDeclare(
		rl.config.RabbitMQQueue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{"x-dead-letter-exchange": rl.config.DLXExchange},
	)
	if err != nil {
		return fmt.Errorf("failed to declare raw events queue: %w", err)
	}
	log.Printf("Declared raw events queue: %s (DLX: %s)", q.Name, rl.config.DLXExchange)

	if err = ch.QueueBind(q.Name, "", rl.config.RawEventsExchange, false, nil); err != nil {
		return fmt.Errorf("failed to bind raw events queue to exchange: %w", err)
	}
	log.Printf("Bound queue %s to exchange %s", q.Name, rl.config.RawEventsExchange)

	// Declare processed events fanout exchange
	if err = ch.ExchangeDeclare(rl.config.ProcessedEventsExchange, "fanout", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare processed events exchange: %w", err)
	}
	log.Printf("Declared processed events exchange: %s", rl.config.ProcessedEventsExchange)

	rl.channel = ch
	rl.queue = &q
	return nil
}
