package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	"github.com/hashicorp/go-multierror"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = ""
	contentType  = "application/octet-stream"
)

type connector struct {
	dsn        string
	queueName  string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func newConnector(dsn, queueName string) *connector {
	return &connector{dsn: dsn, queueName: queueName}
}

func (c *connector) connect() error {
	if c.connection == nil {
		connection, err := amqp.Dial(c.dsn)
		if err != nil {
			return err
		}
		c.connection = connection
	}
	if c.channel == nil {
		ch, err := c.connection.Channel()
		if err != nil {
			return err
		}
		_, err = ch.QueueDeclare(c.queueName, false, false, false, false, nil)
		if err != nil {
			return err
		}
		c.channel = ch
	}
	return nil
}

func (c *connector) close() error {
	var err error
	if c.channel != nil {
		errClose := c.channel.Close()
		if errClose != nil {
			err = multierror.Append(err, errClose)
		}
	}
	if c.connection != nil {
		errClose := c.connection.Close()
		if errClose != nil {
			err = multierror.Append(err, errClose)
		}
	}
	return err
}

func (c *connector) publish(ctx context.Context, notification app.Event) error {
	if c.channel == nil {
		if err := c.connect(); err != nil {
			return fmt.Errorf("connection error: %w", err)
		}
	}
	body, err := eventToBytes(notification)
	if err != nil {
		return fmt.Errorf("notification encode error: %w", err)
	}
	return c.channel.PublishWithContext(
		ctx,
		exchangeName,
		c.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)
}

func eventToBytes(event app.Event) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(event); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *connector) consume(ctx context.Context) (<-chan app.Event, error) {
	if c.channel == nil {
		if err := c.connect(); err != nil {
			return nil, fmt.Errorf("connection error: %w", err)
		}
	}
	ch := make(chan app.Event)
	deliveries, err := c.channel.Consume(c.queueName, "", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("start consuming: %w", err)
	}
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			select {
			case <-ctx.Done():
				_ = c.close()
				return
			default:
				select {
				case <-ctx.Done():
					return
				case delivery := <-deliveries:
					event, err := bytesToEvent(delivery.Body)
					if err != nil {
						// TODO: handle
					} else {
						ch <- event
					}
				}
			}
		}
	}()
	return ch, nil
}

func bytesToEvent(b []byte) (app.Event, error) {
	buffer := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buffer)
	var n app.Event
	if err := decoder.Decode(&n); err != nil {
		return n, fmt.Errorf("notification decode error: %w", err)
	}
	return n, nil
}
