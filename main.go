package main

import (
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config/connections"
)

func main() {
	// Parse environment variables
	config.GetEnvironment()

	// Connect to database
	connections.GetPostgresConnection()
	defer connections.ClosePostgresConnection()

	// Connect to RabbitMQ
	connections.GetRabbitMQChannel()
	defer connections.CloseRabbitMQConnection()
}
