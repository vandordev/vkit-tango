package river

import (
	"context"
	"testing"

	riverqueue "github.com/riverqueue/river"
	"github.com/vandordev/vkit-tango/internal/contract"
	platformrealtime "github.com/vandordev/vkit-tango/internal/platform/realtime"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

func TestRegisterWorkersInstallsRealtimePublisher(t *testing.T) {
	workers := riverqueue.NewWorkers()
	NewRealtimePublishRegistrar(platformrealtime.HTTPPublisher{}).RegisterWorkers(workers)
	if workers == nil {
		t.Fatal("workers = nil")
	}
}

var _ contract.WorkerRegistrar = (*SetSystemMetadataRegistrar)(nil)
var _ contract.WorkerRegistrar = (*RealtimePublishRegistrar)(nil)

func TestSetSystemMetadataWorkerInvokesSharedUsecase(t *testing.T) {
	called := false
	worker := SetSystemMetadataWorker{Command: metadataMutation(func(_ context.Context, input usecase.SetSystemMetadataInput) (usecase.SetSystemMetadataResult, error) {
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
