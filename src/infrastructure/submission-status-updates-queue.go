package infrastructure

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/application"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config/connections"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/dtos"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/infrastructure/implementations"
)

type SubmissionStatusUpdatesQueueMgr struct {
	Queue    *amqp.Queue
	Channel  <-chan amqp.Delivery
	UseCases *application.SubmissionsStatusUpdaterUseCases
}

// submissionStatusUpdatesQueueMgrInstance singleton instance of SubmissionStatusUpdatesQueueMgr
var submissionStatusUpdatesQueueMgrInstance *SubmissionStatusUpdatesQueueMgr

// GetSubmissionStatusUpdatesQueueMgr returns a singleton instance of SubmissionStatusUpdatesQueueMgr
func GetSubmissionStatusUpdatesQueueMgr() *SubmissionStatusUpdatesQueueMgr {
	if submissionStatusUpdatesQueueMgrInstance == nil {
		submissionStatusUpdatesQueueMgrInstance = &SubmissionStatusUpdatesQueueMgr{
			Queue: getRabbitMQSubmissionsStatusUpdatesQueue(),
			UseCases: &application.SubmissionsStatusUpdaterUseCases{
				SubmissionsRepository:           implementations.GetSubmissionsPgRepository(),
				SubmissionsRealTimeUpdatesQueue: implementations.GetSubmissionsRealTimeUpdatesQueueMgrInstance(),
			},
			// Channel is set when listening for submission status updates
		}
	}

	return submissionStatusUpdatesQueueMgrInstance
}

// getRabbitMQSubmissionsStatusUpdatesQueue returns a queue for submission status updates
func getRabbitMQSubmissionsStatusUpdatesQueue() *amqp.Queue {
	ch := connections.GetRabbitMQChannel()

	// Declare queue
	qName := "submission-status-updates"
	qDurable := true
	qAutoDelete := false
	qExclusive := false
	qNoWait := false
	qArgs := amqp.Table{}

	q, err := ch.QueueDeclare(
		qName,
		qDurable,
		qAutoDelete,
		qExclusive,
		qNoWait,
		qArgs,
	)

	if err != nil {
		log.Fatal(
			"[RabbitMQ]: Error declaring submissions status updates queue: ",
			err.Error(),
		)
	}

	// Set fair dispatch
	maxPrefetchCount := 8 // Limit to 8 updates per worker
	err = ch.Qos(
		maxPrefetchCount,
		0,
		false,
	)

	if err != nil {
		log.Fatal(
			"[RabbitMQ]: Error setting fair dispatch for submissions status updates queue: ",
			err.Error(),
		)
	}

	// Set queue
	log.Println("[RabbitMQ]: Submissions status updates queue declared")
	return &q
}

// ListenForSubmissionStatusUpdates starts listening for submission status updates
func (queueMgr *SubmissionStatusUpdatesQueueMgr) ListenForSubmissionStatusUpdates() {
	// Get channel
	ch := connections.GetRabbitMQChannel()

	// Consume
	qName := queueMgr.Queue.Name
	qConsumer := ""
	qAutoAck := false
	qExclusive := false
	qNoLocal := false
	qNoWait := false
	qArgs := amqp.Table{}

	msgs, err := ch.Consume(
		qName,
		qConsumer,
		qAutoAck,
		qExclusive,
		qNoLocal,
		qNoWait,
		qArgs,
	)
	if err != nil {
		log.Fatal(
			"[RabbitMQ Submissions Status Updates Queue]: Error consuming queue: ",
			err.Error(),
		)
	}

	// Set channel
	queueMgr.Channel = msgs

	// Process
	queueMgr.processSubmissionStatusUpdates()
}

// processSubmissionStatusUpdates creates an infinite loop that processes submission status updates
func (queueMgr *SubmissionStatusUpdatesQueueMgr) processSubmissionStatusUpdates() {
	log.Println(
		"[RabbitMQ Submissions Status Updates Queue]: Listening for submission status updates...",
	)

	// Process each message in a separate goroutine
	for msg := range queueMgr.Channel {
		go queueMgr.processSubmissionStatusUpdate(msg)
	}
}

func (queueMgr *SubmissionStatusUpdatesQueueMgr) processSubmissionStatusUpdate(msg amqp.Delivery) {
	// Acknowledge message after processing
	msg.Ack(false)

	// Parse message
	dto := &dtos.SubmissionStatusUpdateDTO{}
	err := json.Unmarshal(msg.Body, dto)
	if err != nil {
		log.Println(
			"[RabbitMQ Submissions Status Updates Queue]: Error parsing submission status update: ",
			err.Error(),
		)
		return
	}

	// Log message to console
	log.Println(
		"[RabbitMQ]: Received submission status update",
		dto.SubmissionUUID,
	)

	// Update submission status
	err = queueMgr.UseCases.UpdateSubmissionStatus(dto)
	if err != nil {
		return
	}
}
