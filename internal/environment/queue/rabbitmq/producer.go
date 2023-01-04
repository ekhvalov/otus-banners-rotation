package rabbitmq

import (
	"context"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
)

func NewProducer(config *Config) *Producer {
	return &Producer{connector: newConnector(config.GetDSN(), config.GetQueueName())}
}

type Producer struct {
	connector *connector
}

func (p *Producer) Put(ctx context.Context, event app.Event) error {
	return p.connector.publish(ctx, event)
}

func (p *Producer) Close() error {
	return p.connector.close()
}
