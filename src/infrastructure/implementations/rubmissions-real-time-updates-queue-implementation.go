package implementations

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config/connections"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/dtos"
)

type SubmissionsRealTimeUpdatesQueueMgr struct {
	Queue *amqp.Queue
}

// ## Singleton instance ##

// submissionsRealTimeUpdatesQueueMgr Singleton instance
var submissionsRealTimeUpdatesQueueMgr *SubmissionsRealTimeUpdatesQueueMgr

// GetSubmissionsRealTimeUpdatesQueueMgrInstance returns the singleton instance of the submissionsRealTimeUpdatesQueueMgr
func GetSubmissionsRealTimeUpdatesQueueMgrInstance() *SubmissionsRealTimeUpdatesQueueMgr {
	if submissionsRealTimeUpdatesQueueMgr == nil {
		submissionsRealTimeUpdatesQueueMgr = &SubmissionsRealTimeUpdatesQueueMgr{
			Queue: getSubmissionsRealTimeUpdatesQueue(),
		}
	}

	return submissionsRealTimeUpdatesQueueMgr
}

// getSubmissionsRealTimeUpdatesQueue returns a pointer to the queue
func getSubmissionsRealTimeUpdatesQueue() *amqp.Queue {
	noQueueHasBeenDeclared :=
		submissionsRealTimeUpdatesQueueMgr == nil ||
			submissionsRealTimeUpdatesQueueMgr.Queue == nil

	if noQueueHasBeenDeclared {
		ch := connections.GetRabbitMQChannel()

		// Declare queue
		qName := "submission-real-time-updates"
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
			log.Fatal(err.Error())
		}

		// Set fair dispatch
		maxPrefetchCount := 4 // Limit to 4 updates per worker
		err = ch.Qos(
			maxPrefetchCount,
			0,
			false,
		)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("[RabbitMQ]: Submission real time updates queue declared")
		return &q
	}

	return submissionsRealTimeUpdatesQueueMgr.Queue
}

// ## Methods implementation ##
// EnqueueUpdate enqueues a submission status update
func (queueMgr *SubmissionsRealTimeUpdatesQueueMgr) EnqueueUpdate(updateDTO *dtos.SubmissionStatusUpdateDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := connections.GetRabbitMQChannel()

	// Parse update to JSON
	json, err := json.Marshal(updateDTO)
	if err != nil {
		log.Println(
			"[RabbitMQ Submissions Real Time Updates]: Error parsing submission status update to JSON",
			err.Error(),
		)
		return err
	}

	// Publish
	mshExchange := ""
	mshRoutingKey := queueMgr.Queue.Name
	mshMandatory := false
	mshImmediate := false

	err = ch.PublishWithContext(
		ctx,
		mshExchange,
		mshRoutingKey,
		mshMandatory,
		mshImmediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        json,
		},
	)

	if err != nil {
		log.Println(
			"[RabbitMQ Submissions Real Time Updates]: Error publishing submission status update",
			err.Error(),
		)
		return err
	}

	return nil
}
