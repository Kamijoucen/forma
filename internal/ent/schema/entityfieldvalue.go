package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// EntityFieldValue holds the schema definition for the EntityFieldValue entity.
type EntityFieldValue struct {
	ent.Schema
}

// Fields of the EntityFieldValue.
func (EntityFieldValue) Fields() []ent.Field {
	return []ent.Field{
		field.Text("value").Default("").Comment("字段值的字符串表示"),
	}
}

// Edges of the EntityFieldValue.
func (EntityFieldValue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("entityRecord", EntityRecord.Type).Ref("fieldValues").Unique().Required().Comment("所属实体记录"),
		edge.From("fieldDef", FieldDef.Type).Ref("fieldValues").Unique().Required().Comment("关联字段定义"),
	}
}

// Mixin of the EntityFieldValue.
func (EntityFieldValue) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
