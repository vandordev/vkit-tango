package river

import (
	"context"

	riverqueue "github.com/riverqueue/river"
	platformrealtime "github.com/vandordev/vkit-fast/internal/platform/realtime"
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

func NewWorkers(publisher platformrealtime.Publisher) *riverqueue.Workers {
	workers := riverqueue.NewWorkers()
	riverqueue.AddWorker(workers, &RealtimePublishWorker{Publisher: publisher})
	return workers
}
