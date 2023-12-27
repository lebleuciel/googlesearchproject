package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MinLen(1).
			MaxLen(512),
		field.Int("user_id"),
		field.String("uuid").
			NotEmpty().
			Unique().
			MinLen(1).
			MaxLen(64),
		field.Int("size"),
		field.String("type").
			NotEmpty().
			MinLen(1).
			MaxLen(512),
		field.Time("created_at").
			Default(time.Now).
			Optional().
			Nillable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deleted_at").
			Default(time.Now).
			Optional().
			Nillable(),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Field("user_id").
			Ref("files").
			Unique().
			Required().
			StructTag(`json:"files"`),
		edge.From("filetype", Filetype.Type).
			Field("type").
			Ref("files").
			Unique().
			Required().
			StructTag(`json:"files"`),

		edge.To("tags", Tag.Type),
	}
}
