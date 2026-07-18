package river

import (
	"context"

	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

type SetSystemMetadataArgs struct {
	Key   string         `json:"key"`
	Value map[string]any `json:"value"`
}

func (SetSystemMetadataArgs) Kind() string { return "system_metadata.set.v1" }

type SetSystemMetadataWorker struct {
	riverqueue.WorkerDefaults[SetSystemMetadataArgs]
	Mutation usecase.SetSystemMetadata
}

func (worker SetSystemMetadataWorker) Work(ctx context.Context, job *riverqueue.Job[SetSystemMetadataArgs]) error {
	_, err := worker.Mutation.Execute(ctx, usecase.SetSystemMetadataInput{Key: job.Args.Key, Value: job.Args.Value})
	return err
}
