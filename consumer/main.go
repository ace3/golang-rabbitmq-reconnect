package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	godotenv.Load(".env")
	connectToRabbitMQ()
}
func connectToRabbitMQ() {
	connection, err := amqp.Dial(os.Getenv("AMQP_DSN"))
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ, retrying...")
		time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
		connectToRabbitMQ()
		return
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// Declare a queue
	queue, err := channel.QueueDeclare(
		"testing", // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Printf("Failed to declare a queue: %s\n", err)
		return
	}

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body)
		}

		// If the range loop exits, the channel has been closed
		fmt.Println("Channel closed, reconnecting...")
		connectToRabbitMQ()
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
