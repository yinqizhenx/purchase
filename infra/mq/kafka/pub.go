package kafka

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/google/wire"
	"github.com/segmentio/kafka-go"

	"purchase/infra/mq"
)

var (
	_ mq.Publisher  = (*kafkaPublisher)(nil)
	_ mq.Subscriber = (*kafkaSubscriber)(nil)
	// _ Event      = (*Message)(nil)
)

const HeaderPropertyKey = "props"

var ProviderSet = wire.NewSet(NewKafkaPublisher, NewKafkaSubscriber)

type kafkaPublisher struct {
	writer *kafka.Writer
	// idg    mq.IDGenFunc
}

func (s *kafkaPublisher) Publish(ctx context.Context, msg *mq.Message) error {
	kmsg, err := s.buildKafkaMessage(msg)
	if err != nil {
		return err
	}
	return s.writer.WriteMessages(ctx, *kmsg)
}

func (s *kafkaPublisher) buildKafkaMessage(m *mq.Message) (*kafka.Message, error) {
	header := make([]kafka.Header, 0)
	header = append(header, kafka.Header{
		Key:   mq.MessageID,
		Value: []byte(m.ID),
	})

	p, err := json.Marshal(m.Header())
	if err != nil {
		return nil, err
	}
	header = append(header, kafka.Header{
		Key:   HeaderPropertyKey,
		Value: p,
	})

	return &kafka.Message{
		Key:     []byte(m.BizCode()),
		Value:   m.Body,
		Topic:   m.EventName(),
		Headers: header,
	}, nil
}

func (s *kafkaPublisher) Close() error {
	return s.writer.Close()
}

func NewKafkaPublisher(c config.Config) (mq.Publisher, error) {
	address := make([]string, 0)
	err := c.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}
	w := &kafka.Writer{
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.LeastBytes{},
	}
	return &kafkaPublisher{writer: w}, nil
}

func SetMessageHeader(m *kafka.Message, k string, v interface{}) error {
	val, err := json.Marshal(v)
	if err != nil {
		return err
	}
	// 存在重复key，直接覆盖
	for i, header := range m.Headers {
		if header.Key == k {
			m.Headers[i] = kafka.Header{
				Key:   k,
				Value: val,
			}
			return nil
		}
	}
	m.Headers = append(m.Headers, kafka.Header{
		Key:   k,
		Value: val,
	})
	return nil
}

// GetMessageHeader 获取消息header值， val需为指针
func GetMessageHeader(m kafka.Message, k string, val interface{}) error {
	for _, header := range m.Headers {
		if header.Key == k {
			return json.Unmarshal(header.Value, val)
		}
	}
	return nil
}
