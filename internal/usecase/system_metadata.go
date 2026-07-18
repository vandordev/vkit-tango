package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vandordev/vkit-fast/internal/platform/db"
	"github.com/vandordev/vkit-fast/internal/platform/db/systemmetadata"
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

type SystemMetadataService struct{ Client *db.Client }

func (service SystemMetadataService) Execute(ctx context.Context, input SetSystemMetadataInput) (SetSystemMetadataResult, error) {
	tx, err := service.Client.Tx(ctx)
	if err != nil {
		return SetSystemMetadataResult{}, err
	}
	defer tx.Rollback()
	entity, err := tx.SystemMetadata.Query().Where(systemmetadata.KeyEQ(input.Key)).Only(ctx)
	if db.IsNotFound(err) {
		entity, err = tx.SystemMetadata.Create().SetKey(input.Key).SetValue(input.Value).Save(ctx)
	} else if err == nil {
		entity, err = tx.SystemMetadata.UpdateOne(entity).SetValue(input.Value).Save(ctx)
	}
	if err != nil {
		return SetSystemMetadataResult{}, err
	}
	if err := tx.Commit(); err != nil {
		return SetSystemMetadataResult{}, err
	}
	return SetSystemMetadataResult{ID: entity.ID, Key: entity.Key, UpdatedAt: entity.UpdatedAt}, nil
}
