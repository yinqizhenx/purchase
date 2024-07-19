package kafka_go

import (
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

// const HeaderPropertyKey = "props"

func NewMessage(m *kafka.Message) *Message {
	if m == nil {
		return nil
	}
	msg := &Message{
		Message: m,
	}
	for _, header := range m.Headers {
		if header.Key == HeaderPropertyKey {
			props := &Property{}
			err := json.Unmarshal(header.Value, props)
			if err == nil {
				msg.props = props
			}
			break
		}
	}
	return msg
}

type Message struct {
	*kafka.Message
	props *Property
}

type Property struct {
	DelayTime      time.Duration `json:"delay_time"`
	RealTopic      string        `json:"real_topic"`
	RetryTopic     string        `json:"retry_topic"`
	ReconsumeTimes int           `json:"reconsume_times"`
	MessageID      string        `json:"message_id"`
	EventName      string        `json:"event_name"`
	RedeliveryTime time.Time     `json:"redelivery_time"`
}

func (m *Message) ToKafkaMessage(topic ...string) (*kafka.Message, error) {
	if len(topic) > 0 {
		m.Message.Topic = topic[0]
	}
	if m.props != nil {
		p, err := json.Marshal(m.props)
		if err != nil {
			return nil, err
		}
		// 存在重复key，直接覆盖
		for i, header := range m.Message.Headers {
			if header.Key == HeaderPropertyKey {
				m.Message.Headers[i] = kafka.Header{
					Key:   HeaderPropertyKey,
					Value: p,
				}
				return m.Message, nil
			}
		}
		m.Message.Headers = append(m.Message.Headers, kafka.Header{
			Key:   HeaderPropertyKey,
			Value: p,
		})
	}
	return m.Message, nil
}

func (m *Message) SetMessageTopic(topic string) {
	m.Message.Topic = topic
}

func (m *Message) SetDelayTime(t time.Duration) {
	if m.props == nil {
		m.props = &Property{
			DelayTime: t,
		}
		return
	}
	m.props.DelayTime = t
}

func (m *Message) SetRealTopic(topic string) {
	if m.props == nil {
		m.props = &Property{
			RealTopic: topic,
		}
		return
	}
	m.props.RealTopic = topic
}

func (m *Message) SetRetryTopic(topic string) {
	if m.props == nil {
		m.props = &Property{
			RetryTopic: topic,
		}
		return
	}
	m.props.RetryTopic = topic
}

func (m *Message) SetReconsumeTimes(n int) {
	if m.props == nil {
		m.props = &Property{
			ReconsumeTimes: n,
		}
		return
	}
	m.props.ReconsumeTimes = n
}

func (m *Message) SetEventName(name string) {
	if m.props == nil {
		m.props = &Property{
			EventName: name,
		}
		return
	}
	m.props.EventName = name
}

func (m *Message) SetMessageID(id string) {
	if m.props == nil {
		m.props = &Property{
			MessageID: id,
		}
		return
	}
	m.props.MessageID = id
}

func (m *Message) SetRedeliveryTime() {
	if m.props == nil {
		m.props = &Property{
			RedeliveryTime: time.Now(),
		}
		return
	}
	m.props.RedeliveryTime = time.Now()
}

func (m *Message) PropsDelayTime() time.Duration {
	if m.props == nil {
		return 0
	}
	return m.props.DelayTime
}

func (m *Message) PropsRealTopic() string {
	if m.props == nil {
		return ""
	}
	return m.props.RealTopic
}

func (m *Message) PropsRetryTopic() string {
	if m.props == nil {
		return ""
	}
	return m.props.RetryTopic
}

func (m *Message) PropsReconsumeTimes() int {
	if m.props == nil {
		return 0
	}
	return m.props.ReconsumeTimes
}

func (m *Message) PropsMessageID() string {
	if m.props == nil {
		return ""
	}
	return m.props.MessageID
}

func (m *Message) PropsRedeliveryTime() time.Time {
	if m.props == nil {
		return time.Time{}
	}
	return m.props.RedeliveryTime
}

func (m *Message) IncrReconsumeTimes() {
	m.props.ReconsumeTimes++
}

func (m *Message) PropsEventName() string {
	if m.props == nil {
		return ""
	}
	return m.props.EventName
}
