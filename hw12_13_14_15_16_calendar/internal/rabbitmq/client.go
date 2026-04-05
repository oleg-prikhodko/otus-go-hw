package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go" //nolint:depguard
)

type QueueClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewQueueClient(addr, username, password, queueName string) (*QueueClient, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s/", username, password, addr)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &QueueClient{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (p *QueueClient) Publish(ctx context.Context, msg any) error {
	return p.PublishTo(ctx, p.queue.Name, msg)
}

func (p *QueueClient) PublishTo(ctx context.Context, queueName string, msg any) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (p *QueueClient) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := p.channel.Consume(
		p.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}
	return msgs, nil
}

func (p *QueueClient) Close() error {
	if err := p.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
