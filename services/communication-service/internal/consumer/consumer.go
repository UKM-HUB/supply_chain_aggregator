// Package consumer wraps the shared RabbitMQ consumer from supply-chain-aggregator/pkg/rabbitmq.
package consumer

import (
	"context"
	"log"

	pkgmq "supply-chain-aggregator/pkg/rabbitmq"
)

// Consumer listens on a RabbitMQ queue.
// Queues are subscribed lazily per Consume call.
type Consumer struct {
	url  string
	noOp bool
}

// NewConsumer stores the broker URL.
// When url is empty it runs in no-op mode (Consume blocks on ctx.Done).
func NewConsumer(url string) (*Consumer, error) {
	if url == "" {
		log.Println("[consumer] RABBITMQ_URL not set, running in no-op mode")
		return &Consumer{noOp: true}, nil
	}
	return &Consumer{url: url}, nil
}

// Consume subscribes to queueName and dispatches messages to handler.
// Returning nil from handler acks the message; any error nacks and requeues it.
// Blocks until ctx is cancelled.
func (c *Consumer) Consume(ctx context.Context, queueName string, handler func(context.Context, []byte) error) error {
	if c.noOp {
		log.Printf("[consumer] no-op mode, skipping queue=%s\n", queueName)
		<-ctx.Done()
		return nil
	}

	inner, err := pkgmq.NewConsumer(c.url, queueName)
	if err != nil {
		return err
	}
	defer inner.Close()

	return inner.Run(ctx, pkgmq.Handler(handler))
}

func (c *Consumer) Close() {}

