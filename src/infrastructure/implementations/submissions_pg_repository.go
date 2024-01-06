package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/upb-code-labs/submissions-status-updater-microservice/src/config/connections"
	"github.com/upb-code-labs/submissions-status-updater-microservice/src/domain/dtos"
)

type SubmissionsPgRepository struct {
	Connection *sql.DB
}

// submissionsPgRepositoryInstance is a singleton of SubmissionsPgRepository
var submissionsPgRepositoryInstance *SubmissionsPgRepository

// GetSubmissionsPgRepository returns a singleton of SubmissionsPgRepository
func GetSubmissionsPgRepository() *SubmissionsPgRepository {
	if submissionsPgRepositoryInstance == nil {
		submissionsPgRepositoryInstance = &SubmissionsPgRepository{
			Connection: connections.GetPostgresConnection(),
		}
	}

	return submissionsPgRepositoryInstance
}

// UpdateSubmissionStatus updates the status of a submission
func (repo *SubmissionsPgRepository) UpdateSubmissionStatus(dto *dtos.SubmissionStatusUpdateDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		UPDATE submissions
		SET status = $1, passing = $2, stdout = $3
		WHERE id = $4
	`

	_, err := repo.Connection.ExecContext(
		ctx,
		query,
		dto.SubmissionStatus,
		dto.TestsPassed,
		dto.TestsOutput,
		dto.SubmissionUUID,
	)

	if err != nil {
		return err
	}

	return nil
}
