package main

import (
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config/connections"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/infrastructure"
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

	// Listen for submission status updates
	submissionStatusUpdatesQueueMgr := infrastructure.GetSubmissionStatusUpdatesQueueMgr()
	submissionStatusUpdatesQueueMgr.ListenForSubmissionStatusUpdates()
}
