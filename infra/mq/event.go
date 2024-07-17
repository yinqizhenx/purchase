package mq

import (
	"context"

	"github.com/google/uuid"
	"github.com/google/wire"
)

// var ProviderSet = wire.NewSet(kafka.NewKafkaPublisher, NewPulsarSubscriber, NewIDGenFunc)
var ProviderSet = wire.NewSet(NewIDGenFunc)

type Publisher interface {
	Publish(context.Context, *Message) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context)
	Close() error
}

type Handler func(context.Context, *Message) error

type Message struct {
	// id     string
	header map[string]string
	Body   []byte
}

func (m *Message) HeaderGet(k string) string {
	if len(m.header) > 0 {
		return m.header[k]
	}
	return ""
}

func (m *Message) HeaderSet(k, v string) {
	if m.header == nil {
		m.header = map[string]string{
			k: v,
		}
	} else {
		m.header[k] = v
	}
}

func (m *Message) Header() map[string]string {
	return m.header
}

// func (m *Message) ID() string {
// 	return m.id
// }
//
// func (m *Message) SetID(id string) {
// 	m.id = id
// }

// func (m *Message) SetBody(body []byte) {
// 	m.body = body
// }

const (
	MessageID = "message_id"
	EntityID  = "entity_id"
	EventName = "event_name"
)

type IDGenFunc func() string

func NewID() string {
	return uuid.New().String()
}

func NewIDGenFunc() IDGenFunc {
	return NewID
}
