package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// SchemaDef holds the schema definition for the SchemaDef entity.
type SchemaDef struct {
	ent.Schema
}

// Fields of the SchemaDef.
func (SchemaDef) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Immutable().Comment("Schema名称"),
		field.String("description").Optional().Comment("Schema描述"),
	}
}

// Edges of the SchemaDef.
func (SchemaDef) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("app", App.Type).Ref("schemaDefs").Unique().Required().Comment("所属App"),
		edge.To("fieldDefs", FieldDef.Type).Comment("Schema包含的字段定义"),
		edge.To("entityRecords", EntityRecord.Type).Comment("Schema下的实体记录"),
	}
}

// Mixin of the SchemaDef.
func (SchemaDef) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Indexes of the SchemaDef.
func (SchemaDef) Indexes() []ent.Index {
	return []ent.Index{
		// 同一 App 下 schema name 唯一
		index.Fields("name").Edges("app").Unique(),
	}
}
