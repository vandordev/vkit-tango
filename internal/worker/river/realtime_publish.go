package river

import (
	"context"

	riverqueue "github.com/riverqueue/river"
	platformrealtime "github.com/vandordev/vkit-fast/internal/platform/realtime"
	"github.com/vandordev/vkit-fast/internal/usecase"
)

type RealtimePublishArgs struct {
	Event platformrealtime.Event `json:"event"`
}

func (RealtimePublishArgs) Kind() string { return "realtime.publish.v1" }

type RealtimePublishWorker struct {
	riverqueue.WorkerDefaults[RealtimePublishArgs]
	Publisher platformrealtime.Publisher
}

func (worker RealtimePublishWorker) Work(ctx context.Context, job *riverqueue.Job[RealtimePublishArgs]) error {
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
