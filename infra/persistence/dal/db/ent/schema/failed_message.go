package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type FailedMessage struct {
	ent.Schema
}

func (FailedMessage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_failed_message"},
	}
}

func (FailedMessage) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("message_id").Comment("消息唯一ID"),
		field.String("topic").Comment("原始topic"),
		field.String("biz_code").Default(""),
		field.String("body").Comment("消息体"),
		field.String("header").Comment("消息header JSON"),
		field.String("error_msg").Default("").Comment("失败原因"),
		field.String("state").Default("pending").Comment("pending/recovered"),
		field.Int("retry_count").Default(0),
		field.Time("created_at").
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
		field.Time("updated_at").
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.MySQL: "datetime",
			}),
	}
}
