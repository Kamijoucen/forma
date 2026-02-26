package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// App holds the schema definition for the App entity.
type App struct {
	ent.Schema
}

// Fields of the App.
func (App) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").NotEmpty().Immutable().Unique().Comment("App唯一标识"),
		field.String("name").NotEmpty().Comment("App显示名称"),
		field.String("description").Optional().Comment("App描述"),
	}
}

// Edges of the App.
func (App) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("schemaDefs", SchemaDef.Type).Comment("App下的Schema定义"),
	}
}

// Mixin of the App.
func (App) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
