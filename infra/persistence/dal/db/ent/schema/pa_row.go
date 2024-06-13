package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// PARow holds the schema definition for the Comment entity.
type PARow struct {
	ent.Schema
}

// Fields of the Comment.
func (PARow) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("doc_code"),
		field.String("order_code"),
		field.String("row_currency"),
		field.String("tax_ratio"),
		field.String("initial_amount"),
		field.Time("created_at").
			Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
		field.Time("updated_at").
			Default(time.Now).SchemaType(map[string]string{
			dialect.MySQL: "datetime",
		}),
	}
}
