package queue

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type Config struct {
	URL string
}

type Message struct {
	Type    string
	Payload interface{}
}

func NewClient(config Config) (*Client, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	client := &Client{
		conn:    conn,
		channel: ch,
	}

	// Declare default exchanges and queues
	if err := client.setupQueues(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

func (c *Client) setupQueues() error {
	// Declare exchanges
	if err := c.channel.ExchangeDeclare(
		"ecommerce", // exchange name
		"topic",     // exchange type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queues
	queues := []string{"orders", "inventory", "notifications"}
	for _, q := range queues {
		if _, err := c.channel.QueueDeclare(
			q,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q, err)
		}

		// Bind queue to exchange
		if err := c.channel.QueueBind(
			q,           // queue name
			q,           // routing key
			"ecommerce", // exchange
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", q, err)
		}
	}

	return nil
}

func (c *Client) Publish(queue string, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := c.channel.Publish(
		"ecommerce", // exchange
		queue,       // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

type ConsumerFunc func(Message) error

func (c *Client) Consume(queue string, handler ConsumerFunc) error {
	msgs, err := c.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for delivery := range msgs {
			var msg Message
			if err := json.Unmarshal(delivery.Body, &msg); err != nil {
				// Log error and nack message
				delivery.Nack(false, true)
				continue
			}

			if err := handler(msg); err != nil {
				// Log error and nack message
				delivery.Nack(false, true)
				continue
			}

			delivery.Ack(false)
		}
	}()

	return nil
}

func (c *Client) Close() error {
	if err := c.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (c *Client) ReturnConnection() *amqp.Connection {
	return c.conn
}
