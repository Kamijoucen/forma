package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// SchemaDef holds the schema definition for the SchemaDef entity.
type SchemaDef struct {
	ent.Schema
}

// Fields of the SchemaDef.
func (SchemaDef) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("Schema名称"),
		field.String("description").Optional().Comment("Schema描述"),
	}
}

// Edges of the SchemaDef.
func (SchemaDef) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("fieldDefs", FieldDef.Type).Comment("Schema包含的字段定义"),
	}
}

// Mixin of the SchemaDef.
func (SchemaDef) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
