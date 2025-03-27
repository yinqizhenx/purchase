package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// Trans holds the schema definition for the Trans entity.
type Trans struct {
	ent.Schema
}

// Fields of the Trans.
func (Trans) Fields() []ent.Field {
	return []ent.Field{
		field.String("state"),
		field.String("execute_state"),
		field.String("name"),
		field.Time("finished_at").
			Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.Time("created_at").
			Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.Time("updated_at").
			Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.String("updated_by"),
		field.String("created_by"),
	}
}

// Edges of the Trans.
func (Trans) Edges() []ent.Edge {
	return nil
}
