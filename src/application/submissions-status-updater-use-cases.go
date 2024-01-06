package application

import (
	"log"

	"github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/definitions"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/dtos"
)

type SubmissionsStatusUpdaterUseCases struct {
	SubmissionsRepository           definitions.SubmissionsRepository
	SubmissionsRealTimeUpdatesQueue definitions.SubmissionsRealTimeUpdatesQueue
}

func (useCases *SubmissionsStatusUpdaterUseCases) UpdateSubmissionStatus(dto *dtos.SubmissionStatusUpdateDTO) error {
	// Update submission status in the database
	err := useCases.SubmissionsRepository.UpdateSubmissionStatus(dto)
	if err != nil {
		log.Println("[SubmissionsStatusUpdaterUseCases] Error updating submission status: " + err.Error())
		return err
	}

	// Send submission status update to the submissions real time updates queue
	err = useCases.SubmissionsRealTimeUpdatesQueue.SendUpdate(dto)
	if err != nil {
		log.Println("[SubmissionsStatusUpdaterUseCases] Error enqueuing submission status update: " + err.Error())
		return err
	}

	return nil
}
