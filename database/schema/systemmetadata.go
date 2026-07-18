package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// SystemMetadata stores platform metadata without representing product domain state.
type SystemMetadata struct {
	ent.Schema
}

func (SystemMetadata) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.String("key").NotEmpty().Unique(),
		field.JSON("value", map[string]any{}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}
