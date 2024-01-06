package definitions

import "github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/dtos"

type SubmissionsRealTimeUpdatesQueue interface {
	// Publishes a message to the queue
	SendUpdate(updateDTO *dtos.SubmissionStatusUpdateDTO) error
}
