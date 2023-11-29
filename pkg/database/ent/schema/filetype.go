package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Filetype holds the schema definition for the Filetype entity.
type Filetype struct {
	ent.Schema
}

// Fields of the Filetype.
func (Filetype) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			StorageKey("type").
			NotEmpty().
			MinLen(1).
			MaxLen(256),
		field.Int("allowed_size").
			Default(10000000),
		field.Bool("is_banned").
			Default(false),
		field.Time("created_at").
			Default(time.Now).
			Optional().
			Nillable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Filetype.
func (Filetype) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
	}
}
