package consumer

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	url     string
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(url string) (*Consumer, error) {
	if url == "" {
		log.Println("[consumer] RABBITMQ_URL not set, running in no-op mode")
		return &Consumer{url: url}, nil
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	log.Println("[consumer] connected to RabbitMQ")
	return &Consumer{url: url, conn: conn, channel: ch}, nil
}

// Consume starts consuming messages from the given queue. The handler is called
// for each message; returning nil acks the message, returning an error nacks it.
// Blocks until ctx is cancelled.
func (c *Consumer) Consume(ctx context.Context, queueName string, handler func(context.Context, []byte) error) error {
	if c.channel == nil {
		log.Printf("[consumer] no RabbitMQ connection, skipping queue=%s\n", queueName)
		<-ctx.Done()
		return nil
	}

	_, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register consumer for queue %s: %w", queueName, err)
	}

	log.Printf("[consumer] listening on queue=%s\n", queueName)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("consumer channel closed for queue %s", queueName)
			}

			if err := handler(ctx, msg.Body); err != nil {
				log.Printf("[consumer] handler error for queue=%s: %v\n", queueName, err)
				msg.Nack(false, false)
				continue
			}

			msg.Ack(false)
		}
	}
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
