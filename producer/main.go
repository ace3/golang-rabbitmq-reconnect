package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error loading .env file")
	}
	fmt.Println("RabbitMQ in Golang: Getting started tutorial")

	connection, err := connectToRabbitMQ(os.Getenv("AMQP_DSN"))
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"testing", // name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}

	// Retry sending the message
	err = retryPublish(channel, "testing", []byte("Test Message"), 5)
	if err != nil {
		panic(err)
	}

	fmt.Println("Queue status:", queue)
	fmt.Println("Successfully published message")
}

func connectToRabbitMQ(amqpDSN string) (*amqp.Connection, error) {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		conn, err := amqp.Dial(amqpDSN)
		if err != nil {
			fmt.Printf("Failed to connect to RabbitMQ, retrying... (%d/%d)\n", i+1, maxRetries)
			time.Sleep(2 * time.Second)
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts", maxRetries)
}

func retryPublish(channel *amqp.Channel, queueName string, body []byte, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		err := channel.PublishWithContext(
			context.Background(),
			"",        // exchange
			queueName, // key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			},
		)
		if err != nil {
			fmt.Printf("Failed to publish message, retrying... (%d/%d)\n", i+1, maxRetries)
			time.Sleep(2 * time.Second)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to publish message after %d attempts", maxRetries)
}
