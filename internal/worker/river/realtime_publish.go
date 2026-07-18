package river

import (
	"context"

	riverqueue "github.com/riverqueue/river"
	platformrealtime "github.com/vandordev/vkit-tango/internal/platform/realtime"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

type RealtimePublishArgs = platformrealtime.PublishArgs

type RealtimePublishWorker struct {
	riverqueue.WorkerDefaults[platformrealtime.PublishArgs]
	Publisher platformrealtime.Publisher
}

func (worker RealtimePublishWorker) Work(ctx context.Context, job *riverqueue.Job[platformrealtime.PublishArgs]) error {
	return worker.Publisher.Publish(ctx, job.Args.Event)
}

// RegisterWorkers installs every typed River worker supported by this runtime.
// Job handlers may read Ent directly, but database mutations must invoke an
// internal/usecase command rather than duplicate product rules.
func RegisterWorkers(publisher platformrealtime.Publisher, mutations ...usecase.SetSystemMetadata) (*riverqueue.Workers, error) {
	workers := riverqueue.NewWorkers()
	riverqueue.AddWorker(workers, &RealtimePublishWorker{Publisher: publisher})
	if len(mutations) == 1 && mutations[0] != nil {
		riverqueue.AddWorker(workers, &SetSystemMetadataWorker{Mutation: mutations[0]})
	}
	return workers, nil
}
