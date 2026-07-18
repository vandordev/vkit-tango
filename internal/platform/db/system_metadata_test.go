package db_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/vandordev/vkit-fast/internal/platform/db"
)

var _ uuid.UUID = db.SystemMetadata{}.ID

func TestSystemMetadataUsesUUIDPrimaryKey(t *testing.T) {}
