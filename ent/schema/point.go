package schema

import (
	"townpoint_bot/ent/schema/common"

	"townpoint_bot/pkg/uuidtools"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/gofrs/uuid"
)

// Point holds the schema definition for the Point entity.
type Point struct {
	ent.Schema
}

// Fields of the Point.
func (Point) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuidtools.New),
		field.String("name").Unique(),
		field.JSON("pictures", []common.PointPicture{}).Annotations(
			entgql.Skip(entgql.SkipMutationCreateInput | entgql.SkipMutationUpdateInput),
		).Optional(),
		field.JSON("videos", []common.PointVideo{}).Annotations(
			entgql.Skip(entgql.SkipMutationCreateInput | entgql.SkipMutationUpdateInput),
		).Optional(),
		field.String("address").Unique(),
		field.String("phone"),
		field.String("description").Optional().Nillable(),
		field.Enum("status").NamedValues(
			"Created", "CREATED",
			"Verified", "VERIFIED",
		),
		field.UUID("town_id", uuid.UUID{}),
	}
}

func (Point) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("town", Town.Type).Unique().Required().Field("town_id"),
	}
}

func (Point) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
	}
}

func (Point) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
