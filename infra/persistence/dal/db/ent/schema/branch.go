package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// Branch holds the schema definition for the Branch entity.
type Branch struct {
	ent.Schema
}

// Fields of the Branch.
func (Branch) Fields() []ent.Field {
	return []ent.Field{
		field.String("code"),
		field.Int("trans_id"),
		field.String("type"),
		field.String("state"),
		field.String("name"),
		field.String("action"),
		field.String("compensate"),
		field.String("action_payload"),
		field.String("compensate_payload"),
		field.String("action_depend"),
		field.String("compensate_depend"),
		field.Bool("is_dead"),
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

// Edges of the Branch.
func (Branch) Edges() []ent.Edge {
	return nil
}
