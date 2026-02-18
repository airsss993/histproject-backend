package worker

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeProcessArchive = "archive:process"
)

type ProcessArchivePayload struct {
	RequestId int
	ArchiveId string
}

func NewProcessArchiveTask(requestId int, archiveId string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProcessArchivePayload{
		RequestId: requestId,
		ArchiveId: archiveId,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProcessArchive, payload), nil
}
