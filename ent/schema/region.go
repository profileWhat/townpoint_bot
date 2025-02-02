package schema

import (
	"townpoint_bot/pkg/uuidtools"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/gofrs/uuid"
)

// Region holds the schema definition for the Region entity.
type Region struct {
	ent.Schema
}

// Fields of the Region.
func (Region) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuidtools.New),
		field.String("name").Unique(),
		field.String("description").Optional().Nillable(),
		field.Enum("status").NamedValues(
			"Created", "CREATED",
			"Verified", "VERIFIED",
		),
	}
}

func (Region) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
	}
}

func (Region) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("towns", Town.Type).Ref("region"),
	}
}

func (Region) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
