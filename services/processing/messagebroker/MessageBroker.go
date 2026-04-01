package messagebroker

import (
	"errors"

	"github.com/IshaySela/israel-osint-ai/services/processing/config"
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
	config  *config.Config
}

func NewRabbitListener(config *config.Config) RabbitClient {
	return RabbitClient{config: config}
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
		true, // durable
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

func (rl *RabbitClient) Listen(action func(models.RawOsintEvent)) error {

	if err := rl.setup(); err != nil {
		return err
	}

	msgs, err := rl.channel.Consume(
		rl.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
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
			var event models.RawOsintEvent
			err := event.Unmarshal(d.Body)
			if err != nil {
				continue
			}
			action(event)
		}
	}()

	return nil
}

func (rl *RabbitClient) Publish(exchange string, routingKey string, body []byte) error {
	if rl.channel == nil {
		return errors.New("Failed to publish message, the function was called before creating the channel")
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

func (rl RabbitClient) PublishProcessedEvent()
