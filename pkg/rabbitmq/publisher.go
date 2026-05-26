package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher sends messages to a RabbitMQ queue.
// When created with an empty URL it operates in no-op mode, printing to stdout
// instead of connecting — safe for local development without a broker.
type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewPublisher dials RabbitMQ and opens a channel.
// Returns a no-op Publisher (no error) when url is empty.
func NewPublisher(url string) (*Publisher, error) {
	if url == "" {
		return &Publisher{}, nil
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

	return &Publisher{conn: conn, channel: ch}, nil
}

// Publish sends body to queueName, declaring the queue if it does not exist.
// In no-op mode the message is printed to stdout.
func (p *Publisher) Publish(ctx context.Context, queueName string, body []byte) error {
	if p.channel == nil {
		fmt.Printf("[rabbitmq] no-op publish queue=%s body=%s\n", queueName, string(body))
		return nil
	}

	if _, err := p.channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("rabbitmq: declare queue %s: %w", queueName, err)
	}

	return p.channel.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// Close releases the channel and connection.
func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
