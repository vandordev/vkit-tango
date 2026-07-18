package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vandordev/vkit-fast/internal/platform/db"
	"github.com/vandordev/vkit-fast/internal/platform/db/systemmetadata"
	platformrealtime "github.com/vandordev/vkit-fast/internal/platform/realtime"
)

type SetSystemMetadataInput struct {
	Key   string
	Value map[string]any
}
type SetSystemMetadataResult struct {
	ID        uuid.UUID
	Key       string
	UpdatedAt time.Time
}

// SetSystemMetadata is the write-side boundary for the baseline's sole
// platform entity. Product mutations follow this same intent-specific pattern.
type SetSystemMetadata interface {
	Execute(context.Context, SetSystemMetadataInput) (SetSystemMetadataResult, error)
}

type SystemMetadataService struct{ Runner Runner }

func (service SystemMetadataService) Execute(ctx context.Context, input SetSystemMetadataInput) (SetSystemMetadataResult, error) {
	var result SetSystemMetadataResult
	err := service.Runner.WithinTransaction(ctx, func(ctx context.Context, tx Transaction) error {
		entity, err := tx.Ent.SystemMetadata.Query().Where(systemmetadata.KeyEQ(input.Key)).Only(ctx)
		if db.IsNotFound(err) {
			entity, err = tx.Ent.SystemMetadata.Create().SetKey(input.Key).SetValue(input.Value).Save(ctx)
		} else if err == nil {
			entity, err = tx.Ent.SystemMetadata.UpdateOne(entity).SetValue(input.Value).Save(ctx)
		}
		if err != nil {
			return err
		}
		if _, err := service.Runner.River.InsertTx(ctx, tx.SQL, platformrealtime.PublishArgs{Event: platformrealtime.Event{Type: platformrealtime.ResourceUpdatedV1, EventID: uuid.NewString(), OccurredAt: time.Now().UTC().Format(time.RFC3339Nano), ResourceID: entity.ID.String(), WorkspaceID: "system"}}, nil); err != nil {
			return err
		}
		result = SetSystemMetadataResult{ID: entity.ID, Key: entity.Key, UpdatedAt: entity.UpdatedAt}
		return nil
	})
	return result, err
}
