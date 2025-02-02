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

// Town holds the schema definition for the Town entity.
type Town struct {
	ent.Schema
}

// Fields of the Town.
func (Town) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuidtools.New),
		field.String("name").Unique(),
		field.String("description").Optional().Nillable(),
		field.Enum("status").NamedValues(
			"Created", "CREATED",
			"Verified", "VERIFIED",
		),
		field.UUID("region_id", uuid.UUID{}),
	}
}

func (Town) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("region", Region.Type).Unique().Required().Field("region_id"),
		edge.From("points", Point.Type).Ref("town"),
	}
}

func (Town) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
	}
}

func (Town) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
