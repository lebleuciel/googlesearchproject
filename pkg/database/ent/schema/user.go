package schema

import (
	"errors"
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/lebleuciel/maani/models"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name").
			NotEmpty().
			MinLen(2).
			MaxLen(64),
		field.String("last_name").
			NotEmpty().
			MinLen(2).
			MaxLen(64),
		field.String("email").
			Optional().
			Nillable().
			Unique().
			Validate(func(s string) error {
				emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
				isValid := emailRegex.MatchString(s)
				if !isValid {
					return errors.New("email is not valid")
				}
				return nil
			}),
		field.String("password").NotEmpty(),
		field.Enum("access_type").
			Values(models.AdminType, models.CustomerType).
			Default(models.CustomerType),
		field.Time("created_at").
			Default(time.Now).
			Optional().
			Nillable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("last_login_at").
			Optional().
			Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
	}
}
