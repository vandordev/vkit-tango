package river

import (
	"context"

	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/contract"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

type SetSystemMetadataArgs struct {
	Key   string         `json:"key"`
	Value map[string]any `json:"value"`
}

func (SetSystemMetadataArgs) Kind() string { return "system_metadata.set.v1" }

type SetSystemMetadataWorker struct {
	riverqueue.WorkerDefaults[SetSystemMetadataArgs]
	Command contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult]
}

func (worker SetSystemMetadataWorker) Work(ctx context.Context, job *riverqueue.Job[SetSystemMetadataArgs]) error {
	_, err := worker.Command.Execute(ctx, usecase.SetSystemMetadataInput{Key: job.Args.Key, Value: job.Args.Value})
	return err
}

type SetSystemMetadataRegistrar struct {
	command contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult]
}

var _ contract.WorkerRegistrar = (*SetSystemMetadataRegistrar)(nil)

func NewSetSystemMetadataRegistrar(command contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult]) *SetSystemMetadataRegistrar {
	return &SetSystemMetadataRegistrar{command: command}
}

func (registrar *SetSystemMetadataRegistrar) RegisterWorkers(workers *riverqueue.Workers) {
	riverqueue.AddWorker(workers, &SetSystemMetadataWorker{Command: registrar.command})
}
