package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Message holds the schema definition for the event.
type AsyncTask struct {
	ent.Schema
}

// Annotations Message实体的注解
func (AsyncTask) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_async_task"},
	}
}

// Fields of the Post.
func (AsyncTask) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("task_id"),
		field.String("task_type"),
		field.String("task_group"),
		field.String("task_name"),
		field.String("biz_id"),
		field.String("task_data"),
		field.String("state").StorageKey("task_state"),
		field.Int("retry_count").Default(0),
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
