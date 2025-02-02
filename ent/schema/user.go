package schema

import (
	"townpoint_bot/pkg/uuidtools"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/gofrs/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuidtools.New).Unique(),
		field.String("tg_id").Unique(),
		field.Enum("role").NamedValues(
			"Admin", "ADMIN",
			"Other", "OTHER",
		),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
