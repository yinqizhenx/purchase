package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Idempotent struct {
	ent.Schema
}

// Annotations Message实体的注解
func (Idempotent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_idempotent"},
	}
}

// Fields of the Post.
func (Idempotent) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("type"),
		field.String("key"),
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

func (Idempotent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type", "key").Unique().StorageKey("idx_type_key"),
	}
}
