package mq

import (
	"context"
	"time"

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

type MessageCommitter interface {
	CommitMessage(context.Context, *Message) error
}

type MessageWriter interface {
	WriteMessages(context.Context, *Message) error
}

type Handler func(context.Context, *Message) error

type Header struct {
	BizCode        string        `json:"biz_code"`
	DelayTime      time.Duration `json:"delay_time"`
	OriginalTopic  string        `json:"original_topic"`
	RetryTopic     string        `json:"retry_topic"`
	DeadTopic      string        `json:"dead_topic"`
	ReconsumeTimes int           `json:"reconsume_times"`
	MessageID      string        `json:"message_id"`
	EventName      string        `json:"event_name"`
	RedeliveryTime time.Time     `json:"redelivery_time"`
}

type Message struct {
	ID         string
	Body       []byte
	header     Header
	rawMessage interface{}
}

// func (m *Message) HeaderGet(k string) interface{} {
// 	if len(m.header) > 0 {
// 		return m.header[k]
// 	}
// 	return nil
// }

// func (m *Message) HeaderSet(k string, v interface{}) {
// 	if m.header == nil {
// 		m.header = map[string]interface{}{
// 			k: v,
// 		}
// 	} else {
// 		m.header[k] = v
// 	}
// }

func (m *Message) Header() Header {
	return m.header
}

func (m *Message) SetHeader(h Header) {
	m.header = h
}

func (m *Message) RawMessage() interface{} {
	return m.header
}

func (m *Message) SetRawMessage(msg interface{}) {
	m.rawMessage = msg
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
	// BizCode   = "biz_code"
	EventName = "event_name"
)

// func (m *Message) SetMessageTopic(topic string) {
// 	m.Message.Topic = topic
// }

func (m *Message) BizCode() string {
	return m.header.BizCode
}

func (m *Message) SetBizCode(code string) {
	m.header.BizCode = code
}

func (m *Message) DelayTime() time.Duration {
	return m.header.DelayTime
}

func (m *Message) SetDelayTime(t time.Duration) {
	m.header.DelayTime = t
}

func (m *Message) OriginalTopic() string {
	return m.header.OriginalTopic
}

func (m *Message) SetOriginalTopic(topic string) {
	m.header.OriginalTopic = topic
}

func (m *Message) RetryTopic() string {
	return m.header.RetryTopic
}

func (m *Message) SetRetryTopic(topic string) {
	m.header.RetryTopic = topic
}

func (m *Message) DeadTopic() string {
	return m.header.DeadTopic
}

func (m *Message) SetDeadTopic(topic string) {
	m.header.DeadTopic = topic
}

func (m *Message) ReconsumeTimes() int {
	return m.header.ReconsumeTimes
}

func (m *Message) SetReconsumeTimes(n int) {
	m.header.ReconsumeTimes = n
}

func (m *Message) EventName() string {
	return m.header.EventName
}

func (m *Message) SetEventName(name string) {
	m.header.EventName = name
}

func (m *Message) RedeliveryTime() time.Time {
	return m.header.RedeliveryTime
}

func (m *Message) SetRedeliveryTime() {
	m.header.RedeliveryTime = time.Now()
}

func (m *Message) IncrReconsumeTimes() {
	m.header.ReconsumeTimes++
}

type IDGenFunc func() string

func NewID() string {
	return uuid.New().String()
}

func NewIDGenFunc() IDGenFunc {
	return NewID
}
