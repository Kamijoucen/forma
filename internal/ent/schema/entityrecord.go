package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/mixin"
)

// EntityRecord holds the schema definition for the EntityRecord entity.
type EntityRecord struct {
	ent.Schema
}

// Fields of the EntityRecord.
func (EntityRecord) Fields() []ent.Field {
	return nil
}

// Edges of the EntityRecord.
func (EntityRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("schemaDef", SchemaDef.Type).Ref("entityRecords").Unique().Required().Comment("所属Schema"),
		edge.To("fieldValues", EntityFieldValue.Type).Comment("实体的字段值"),
	}
}

// Mixin of the EntityRecord.
func (EntityRecord) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
