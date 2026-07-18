package river

import (
	"context"
	"errors"
	"testing"

	riverqueue "github.com/riverqueue/river"
	platformrealtime "github.com/vandordev/vkit-tango/internal/platform/realtime"
)

type publisherFunc func(context.Context, platformrealtime.Event) error

func (fn publisherFunc) Publish(ctx context.Context, event platformrealtime.Event) error {
	return fn(ctx, event)
}

func TestRealtimePublishWorkerReturnsPublisherFailure(t *testing.T) {
	want := errors.New("unavailable")
	worker := RealtimePublishWorker{Publisher: publisherFunc(func(context.Context, platformrealtime.Event) error { return want })}
	job := &riverqueue.Job[RealtimePublishArgs]{Args: RealtimePublishArgs{Event: platformrealtime.Event{Type: platformrealtime.ResourceUpdatedV1}}}
	if err := worker.Work(context.Background(), job); !errors.Is(err, want) {
		t.Fatalf("Work() error = %v, want %v", err, want)
	}
}
