package application

import (
	"log"

	"github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/definitions"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/dtos"
)

type SubmissionsStatusUpdaterUseCases struct {
	SubmissionsRepository definitions.SubmissionsRepository
}

func (useCases *SubmissionsStatusUpdaterUseCases) UpdateSubmissionStatus(dto *dtos.SubmissionStatusUpdateDTO) error {
	// Update submission status in the database
	err := useCases.SubmissionsRepository.UpdateSubmissionStatus(dto)
	if err != nil {
		log.Println("[SubmissionsStatusUpdaterUseCases] Error updating submission status: " + err.Error())
		return err
	}

	// TODO: Send submission status update to the submissions real time updates queue
	return nil
}
