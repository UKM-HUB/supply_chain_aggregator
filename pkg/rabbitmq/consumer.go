package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Handler is a function that processes a single message body.
// Return nil to ack, return any error to nack (message will be requeued).
type Handler func(ctx context.Context, body []byte) error

// Consumer listens on a single RabbitMQ queue and dispatches messages to a Handler.
type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

// NewConsumer dials RabbitMQ and opens a channel subscribed to queueName.
// Returns an error when url is empty — consumers must have a real broker.
func NewConsumer(url, queueName string) (*Consumer, error) {
	if url == "" {
		return nil, fmt.Errorf("rabbitmq: RABBITMQ_URL is required for consumer")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: open channel: %w", err)
	}

	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("rabbitmq: declare queue %s: %w", queueName, err)
	}

	return &Consumer{conn: conn, channel: ch, queueName: queueName}, nil
}

// Run blocks consuming messages until ctx is cancelled.
// Each message is dispatched to handler; on success it is acked, on error it is
// nacked and requeued once.
func (c *Consumer) Run(ctx context.Context, handler Handler) error {
	msgs, err := c.channel.Consume(c.queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("rabbitmq: consume %s: %w", c.queueName, err)
	}

	log.Printf("[rabbitmq] consumer started queue=%s", c.queueName)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("rabbitmq: channel closed for queue %s", c.queueName)
			}
			if err := handler(ctx, msg.Body); err != nil {
				log.Printf("[rabbitmq] handler error queue=%s: %v — nacking", c.queueName, err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}
}

// Close releases the channel and connection.
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
