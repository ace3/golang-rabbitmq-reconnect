package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ace3/golang-rabbitmq-reconnect/lib/rabbitmq"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	fmt.Println("RabbitMQ in Golang: Getting started tutorial")
	consumerFunc := func(body []byte) error {
		log.Printf("Received Message: %s\n", body)
		return nil
	}

	consumer := rabbitmq.NewConsumer("testing", consumerFunc)
	consumer.ConnectAndConsume(os.Getenv("AMQP_DSN"))

	select {} // Keep the application running

}
