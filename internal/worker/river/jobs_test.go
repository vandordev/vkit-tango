package river

import (
	"context"
	"testing"

	riverqueue "github.com/riverqueue/river"
	platformrealtime "github.com/vandordev/vkit-tango/internal/platform/realtime"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

func TestRegisterWorkersInstallsRealtimePublisher(t *testing.T) {
	workers, err := RegisterWorkers(platformrealtime.HTTPPublisher{})
	if err != nil {
		t.Fatalf("RegisterWorkers() error = %v", err)
	}
	if workers == nil {
		t.Fatal("RegisterWorkers() workers = nil")
	}
}

func TestSetSystemMetadataWorkerInvokesSharedUsecase(t *testing.T) {
	called := false
	worker := SetSystemMetadataWorker{Mutation: metadataMutation(func(_ context.Context, input usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
		called = input.Key == "feature"
		return usecase.SetSystemMetadataResult{}, nil
	})}
	if err := worker.Work(context.Background(), &riverqueue.Job[SetSystemMetadataArgs]{Args: SetSystemMetadataArgs{Key: "feature"}}); err != nil || !called {
		t.Fatalf("Work() error=%v called=%v", err, called)
	}
}

type metadataMutation func(context.Context, usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error)

func (mutation metadataMutation) Execute(ctx context.Context, input usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
	return mutation(ctx, input)
}
