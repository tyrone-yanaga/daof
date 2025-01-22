package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Client struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	config       Config
	isConnected  bool
	mu           sync.RWMutex
	closed       chan struct{}
	consumers    map[string]ConsumerFunc
	reconnecting bool
}

type Config struct {
	URL               string
	ReconnectInterval time.Duration
	MaxRetries        int
}

type Message struct {
	Type    string
	Payload interface{}
}

func NewClient(config Config) (*Client, error) {
	if config.ReconnectInterval == 0 {
		config.ReconnectInterval = 5 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 10
	}

	client := &Client{
		config:    config,
		closed:    make(chan struct{}),
		consumers: make(map[string]ConsumerFunc),
	}

	// Try to connect with retries during initialization
	var err error
	for i := 0; i < config.MaxRetries; i++ {
		log.Printf("Initial connection attempt %d/%d to RabbitMQ", i+1, config.MaxRetries)
		err = client.connect()
		if err == nil {
			break
		}
		log.Printf("Failed to connect, retrying in %v: %v", config.ReconnectInterval, err)
		time.Sleep(config.ReconnectInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to establish initial connection after %d attempts: %w",
			config.MaxRetries, err)
	}

	go client.handleReconnect()

	return client, nil
}

func (c *Client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.reconnecting {
		log.Printf("Already reconnecting, skipping connect attempt")
		return nil
	}

	log.Printf("Attempting to connect to RabbitMQ at URL: %s", c.config.URL)

	conn, err := amqp.Dial(c.config.URL)
	if err != nil {
		log.Printf("Connection error details: %+v", err)
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	log.Printf("Successfully established connection, opening channel")

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Printf("Channel error details: %+v", err)
		return fmt.Errorf("failed to open channel: %w", err)
	}

	log.Printf("Successfully opened channel, enabling confirmations")

	// Enable channel-level confirmation
	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		log.Printf("Confirmation mode error details: %+v", err)
		return fmt.Errorf("failed to enable channel confirmations: %w", err)
	}

	c.conn = conn
	c.channel = ch
	c.isConnected = true

	log.Printf("RabbitMQ connection fully established")

	// Set up connection monitoring
	go func() {
		closeErr := <-c.channel.NotifyClose(make(chan *amqp.Error))
		log.Printf("Channel closed, error: %+v", closeErr)
		c.mu.Lock()
		c.isConnected = false
		c.mu.Unlock()
	}()

	// Setup queues after successful connection
	if err := c.setupQueues(); err != nil {
		c.Close()
		log.Printf("Queue setup error: %+v", err)
		return err
	}

	log.Printf("Queue setup completed successfully")

	// Restore consumers
	for queue, handler := range c.consumers {
		if err := c.startConsumer(queue, handler); err != nil {
			log.Printf("Failed to restore consumer for queue %s: %v", queue, err)
		}
	}

	return nil
}

func (c *Client) handleReconnect() {
	for {
		select {
		case <-c.closed:
			return
		default:
			c.mu.RLock()
			connected := c.isConnected
			c.mu.RUnlock()

			if !connected {
				c.mu.Lock()
				c.reconnecting = true
				c.mu.Unlock()

				backoff := time.Second
				for attempt := 1; attempt <= c.config.MaxRetries; attempt++ {
					log.Printf("Attempting to reconnect to RabbitMQ (attempt %d/%d)", attempt, c.config.MaxRetries)

					if err := c.connect(); err != nil {
						log.Printf("Failed to reconnect: %v", err)
						time.Sleep(backoff)
						backoff *= 2 // Exponential backoff
						continue
					}

					c.mu.Lock()
					c.reconnecting = false
					c.mu.Unlock()
					log.Printf("Successfully reconnected to RabbitMQ")
					break
				}
			}
			time.Sleep(c.config.ReconnectInterval)
		}
	}
}

func (c *Client) setupQueues() error {
	// Declare exchanges with retry
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

func (c *Client) Publish(ctx context.Context, queue string, msg Message) error {
	c.mu.RLock()
	if !c.isConnected {
		c.mu.RUnlock()
		return fmt.Errorf("not connected to RabbitMQ")
	}
	c.mu.RUnlock()

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	confirms := c.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	if err := c.channel.Publish(
		"ecommerce", // exchange
		queue,       // routing key
		true,        // mandatory
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			return fmt.Errorf("failed to deliver message to queue")
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (c *Client) startConsumer(queue string, handler ConsumerFunc) error {
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
		for {
			select {
			case <-c.closed:
				return
			case delivery, ok := <-msgs:
				if !ok {
					return
				}

				var msg Message
				if err := json.Unmarshal(delivery.Body, &msg); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					delivery.Nack(false, true)
					continue
				}

				if err := handler(msg); err != nil {
					log.Printf("Failed to handle message: %v", err)
					delivery.Nack(false, true)
					continue
				}

				delivery.Ack(false)
			}
		}
	}()

	return nil
}

type ConsumerFunc func(Message) error

func (c *Client) Consume(queue string, handler ConsumerFunc) error {
	c.mu.Lock()
	c.consumers[queue] = handler
	c.mu.Unlock()

	return c.startConsumer(queue, handler)
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.closed)
	c.isConnected = false

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}

	return nil
}

func (c *Client) ReturnConnection() *amqp.Connection {
	return c.conn
}
