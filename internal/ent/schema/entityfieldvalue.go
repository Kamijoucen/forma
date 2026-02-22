package schema

import (
	"forma/internal/constant"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EntityFieldValue holds the schema definition for the EntityFieldValue entity.
type EntityFieldValue struct {
	ent.Schema
}

// Fields of the EntityFieldValue.
func (EntityFieldValue) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Comment("字段名"),
		field.Enum("type").Values(constant.GetFieldTypes()...).Comment("字段类型"),
		field.Text("value").Default("").Comment("字段值的字符串表示"),
	}
}

// Edges of the EntityFieldValue.
func (EntityFieldValue) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("entityRecord", EntityRecord.Type).Ref("fieldValues").Unique().Required().Comment("所属实体记录"),
	}
}
