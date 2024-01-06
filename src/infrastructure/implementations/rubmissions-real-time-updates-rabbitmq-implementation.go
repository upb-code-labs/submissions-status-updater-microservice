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

var exchangeName = "submissions-real-time-updates"

type SubmissionsRealTimeUpdatesRabbitmqMgr struct{}

// ## Singleton instance ##

// submissionsRealTimeUpdatesRabbitmqMgr Singleton instance
var submissionsRealTimeUpdatesRabbitmqMgr *SubmissionsRealTimeUpdatesRabbitmqMgr

// GetSubmissionsRealTimeUpdatesQueueMgrInstance returns the singleton instance of the submissionsRealTimeUpdatesQueueMgr
func GetSubmissionsRealTimeUpdatesQueueMgrInstance() *SubmissionsRealTimeUpdatesRabbitmqMgr {
	if submissionsRealTimeUpdatesRabbitmqMgr == nil {
		// Declare exchange
		declareRealTimeUpdatesExchange()

		// Declare singleton instance
		submissionsRealTimeUpdatesRabbitmqMgr = &SubmissionsRealTimeUpdatesRabbitmqMgr{}
	}

	return submissionsRealTimeUpdatesRabbitmqMgr
}

func declareRealTimeUpdatesExchange() {
	ch := connections.GetRabbitMQChannel()

	exchangeName := exchangeName
	exchangeType := "fanout"
	exchangeDurable := true
	exchangeAutoDelete := false
	exchangeInternal := false
	exchangeNoWait := false
	exchangeArgs := amqp.Table{}

	err := ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		exchangeArgs,
	)

	if err != nil {
		log.Fatal(
			"[RabbitMQ]: Error declaring real time updates exchange",
			err.Error(),
		)
	}

	log.Println("[RabbitMQ]: Real time updates exchange declared")
}

// ## Methods implementation ##

// SendUpdate enqueues a submission status update
func (queueMgr *SubmissionsRealTimeUpdatesRabbitmqMgr) SendUpdate(updateDTO *dtos.SubmissionStatusUpdateDTO) error {
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
	msgExchange := exchangeName
	msgRoutingKey := ""
	msgMandatory := false
	msgImmediate := false

	err = ch.PublishWithContext(
		ctx,
		msgExchange,
		msgRoutingKey,
		msgMandatory,
		msgImmediate,
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
