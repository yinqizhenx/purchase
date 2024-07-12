package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// PAHead holds the schema definition for the Article entity.
type PAHead struct {
	ent.Schema
}

// Fields of the Post.
func (PAHead) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("code").Unique(),
		field.String("state"),
		field.String("pay_amount"),
		field.String("applicant"),
		field.String("department_code").Optional(),
		field.String("supplier_code"),
		field.Bool("is_adv"),
		field.Bool("has_invoice"),
		field.Text("remark"),
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
