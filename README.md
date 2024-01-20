# README
Golang example implemenation of amqp producer & consumer that can persist the restarts of RabbitMQ servers.

## How To Configure
1. Clone the repo
2. Run `go mod download`
3. Run `cp .env.example .env`
4. Configure the .env with dsn `amqp://guest:guest@localhost:5672/` as example
5. Run the Consumer using `go run ./consumer/main.go` , and it will listen for the message
6. Run the Producer using `go run ./producer/main.go` , and it will publish the message to the queue

