package rabbitmq

import (
	"context"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
)

func NewConsumer(config *Config) *Consumer {
	return &Consumer{connector: newConnector(config.GetDSN(), config.GetQueueName())}
}

type Consumer struct {
	connector *connector
}

func (c *Consumer) Subscribe(ctx context.Context) (<-chan app.Event, error) {
	return c.connector.consume(ctx)
}
