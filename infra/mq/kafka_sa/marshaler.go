package kafka_sa

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"

	"purchase/infra/mq"
)

type MessageMarshaler interface {
	Marshal(ctx context.Context, m *mq.Message) (*sarama.ProducerMessage, error)
	Unmarshal(ctx context.Context, m *sarama.ConsumerMessage) (*mq.Message, error)
}

type MarshalerImpl struct {
	idg mq.IDGenFunc
}

func NewMessageMarshaler(idGen mq.IDGenFunc) MessageMarshaler {
	return &MarshalerImpl{}
}

func (ms *MarshalerImpl) Marshal(ctx context.Context, m *mq.Message) (*sarama.ProducerMessage, error) {
	header := make([]sarama.RecordHeader, 0)
	header = append(header, sarama.RecordHeader{
		Key:   []byte(mq.MessageID),
		Value: []byte(ms.idg()),
	})

	p, err := json.Marshal(m.Header())
	if err != nil {
		return nil, err
	}
	header = append(header, sarama.RecordHeader{
		Key:   []byte(HeaderPropertyKey),
		Value: p,
	})

	msg := &sarama.ProducerMessage{
		Key:     sarama.StringEncoder(m.BizCode()),
		Value:   sarama.ByteEncoder(m.Body),
		Topic:   m.EventName(),
		Headers: header,
	}

	// 重投的消息
	if m.RetryTopic() != "" {
		msg.Topic = m.RetryTopic()
	}

	// 死信消息
	if m.DeadTopic() != "" {
		msg.Topic = m.DeadTopic()
	}
	return msg, nil
}

func (ms *MarshalerImpl) Unmarshal(ctx context.Context, m *sarama.ConsumerMessage) (*mq.Message, error) {
	msg := &mq.Message{
		Body: m.Value,
	}
	msg.SetPartition(m.Partition)
	for _, header := range m.Headers {
		if string(header.Key) == mq.MessageID {
			msg.ID = string(header.Value)
		}
		if string(header.Key) == HeaderPropertyKey {
			props := &mq.Header{}
			err := json.Unmarshal(header.Value, props)
			if err != nil {
				return nil, err
			}
			msg.SetHeader(*props)
		}
	}
	msg.SetRawMessage(m)

	return msg, nil
}
