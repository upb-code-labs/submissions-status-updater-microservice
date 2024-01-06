package dtos

type SubmissionStatusUpdateDTO struct {
	SubmissionUUID   string `json:"submission_uuid"`
	SubmissionStatus string `json:"submission_status"`
	TestsPassed      bool   `json:"tests_passed"`
	TestsOutput      string `json:"tests_output"`
}
