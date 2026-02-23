package schema

import (
	"forma/internal/constant"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// FieldDef holds the schema definition for the FieldDef entity.
type FieldDef struct {
	ent.Schema
}

// Fields of the FieldDef.
func (FieldDef) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Immutable().Comment("字段名称"),
		field.Enum("type").Values(constant.GetFieldTypes()...).Immutable().Comment("字段类型"),
		field.Bool("required").Default(false).Comment("是否必填"),
		field.Int("maxLength").Default(500).Comment("最大长度"),
		field.Int("minLength").Default(0).Comment("最小长度"),
		field.JSON("enumValues", []string{}).Optional().Comment("枚举值列表"),
		field.String("description").Optional().Comment("字段描述"),
	}
}

// Edges of the FieldDef.
func (FieldDef) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("schemaDef", SchemaDef.Type).Ref("fieldDefs").Unique().Comment("所属Schema"),
		edge.To("fieldValues", EntityFieldValue.Type).Comment("关联的字段值"),
	}
}

// Mixin of the FieldDef.
func (FieldDef) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
