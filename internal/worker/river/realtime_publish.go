package river

import (
	"context"

	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/contract"
	platformrealtime "github.com/vandordev/vkit-tango/internal/platform/realtime"
)

type RealtimePublishArgs = platformrealtime.PublishArgs

type RealtimePublishWorker struct {
	riverqueue.WorkerDefaults[platformrealtime.PublishArgs]
	Publisher platformrealtime.Publisher
}

func (worker RealtimePublishWorker) Work(ctx context.Context, job *riverqueue.Job[platformrealtime.PublishArgs]) error {
	return worker.Publisher.Publish(ctx, job.Args.Event)
}

type RealtimePublishRegistrar struct{ publisher platformrealtime.Publisher }

var _ contract.WorkerRegistrar = (*RealtimePublishRegistrar)(nil)

func NewRealtimePublishRegistrar(publisher platformrealtime.Publisher) *RealtimePublishRegistrar {
	return &RealtimePublishRegistrar{publisher: publisher}
}

func (registrar *RealtimePublishRegistrar) RegisterWorkers(workers *riverqueue.Workers) {
	riverqueue.AddWorker(workers, &RealtimePublishWorker{Publisher: registrar.publisher})
}
