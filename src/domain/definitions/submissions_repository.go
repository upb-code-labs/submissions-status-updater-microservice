package definitions

import "github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/dtos"

type SubmissionsRepository interface {
	UpdateSubmissionStatus(dto *dtos.SubmissionStatusUpdateDTO) error
}
