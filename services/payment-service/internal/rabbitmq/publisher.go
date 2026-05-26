// Package rabbitmq re-exports the shared Publisher from supply-chain-aggregator/pkg/rabbitmq.
package rabbitmq

import pkgmq "supply-chain-aggregator/pkg/rabbitmq"

// Publisher sends messages to a RabbitMQ queue.
type Publisher = pkgmq.Publisher

// NewPublisher dials RabbitMQ and returns a Publisher.
// Returns a no-op Publisher (no error) when url is empty.
func NewPublisher(url string) (*Publisher, error) {
	return pkgmq.NewPublisher(url)
}
