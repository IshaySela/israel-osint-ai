package messagebroker

import (
	"errors"

	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitListener struct {
	Url       string
	QueueName string
	conn      *amqp.Connection
	channel   *amqp.Channel
	queue     *amqp.Queue
}

func NewRabbitListener(url, queue string) RabbitListener {
	return RabbitListener{Url: url, QueueName: queue}
}

func (rl *RabbitListener) setup() error {
	conn, err := amqp.Dial(rl.Url)
	if err != nil {
		return errors.New("failed to establish connection to rabbitmq host")
	}
	rl.conn = conn
	ch, err := rl.conn.Channel()

	if err != nil {
		return errors.New("failed to open channel to rabbitmq host")
	}

	q, err := ch.QueueDeclare(
		rl.QueueName, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		return errors.New("failed to declare queue")
	}

	rl.channel = ch
	rl.queue = &q
	return nil
}

func (rl *RabbitListener) Listen(action func(models.RawOsintEvent)) error {

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
