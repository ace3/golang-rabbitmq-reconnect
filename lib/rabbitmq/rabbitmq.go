package rabbitmq

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	QueueName  string
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Consumer   func([]byte) error
	amqpURL    string
}

func NewConsumer(queueName string, consumerFunc func([]byte) error) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		QueueName: queueName,
		Consumer:  consumerFunc,
	}
}

func (c *RabbitMQConsumer) ConnectAndConsume(amqpURL string) {
	c.amqpURL = amqpURL
	go c.connect()
}

func (c *RabbitMQConsumer) connect() {
	for {
		conn, err := amqp.Dial(c.amqpURL)
		if err != nil {
			log.Println("Failed to connect to RabbitMQ, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		c.Connection = conn

		ch, err := conn.Channel()
		if err != nil {
			log.Println("Failed to open a channel")
			continue
		}
		c.Channel = ch

		_, err = ch.QueueDeclare(
			c.QueueName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Failed to declare a queue: %s\n", err)
			continue
		}

		msgs, err := ch.Consume(
			c.QueueName,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("Failed to register a consumer: %s\n", err)
			continue
		}

		c.consumeMessages(msgs)
		break
	}
}

func (c *RabbitMQConsumer) consumeMessages(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		err := c.Consumer(msg.Body)
		if err != nil {
			log.Printf("Error processing message: %s", err)
		}
	}

	log.Println("Connection closed. Attempting to reconnect...")
	c.connect()
}
