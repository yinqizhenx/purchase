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

var ProviderSet = wire.NewSet(NewKafkaPublisher, NewKafkaSubscriber)

type kafkaPublisher struct {
	writer *kafka.Writer
	idg    mq.IDGenFunc
}

func (s *kafkaPublisher) Publish(ctx context.Context, msg *mq.Message) error {
	kmsg, err := s.buildKafkaMessage(msg)
	if err != nil {
		return err
	}
	return s.writer.WriteMessages(ctx, *kmsg)
}

func (s *kafkaPublisher) buildKafkaMessage(msg *mq.Message) (*kafka.Message, error) {
	kmsg := &kafka.Message{
		Key:   []byte(msg.HeaderGet(mq.EntityID)),
		Value: msg.Body,
		Topic: msg.HeaderGet(mq.EventName),
	}

	for k, v := range msg.Header() {
		err := SetMessageHeader(kmsg, k, v)
		if err != nil {
			return nil, err
		}
	}
	err := SetMessageHeader(kmsg, mq.MessageID, s.idg())
	if err != nil {
		return nil, err
	}
	return kmsg, nil
}

func (s *kafkaPublisher) Close() error {
	return s.writer.Close()
}

func NewKafkaPublisher(c config.Config, idg mq.IDGenFunc) (mq.Publisher, error) {
	address := make([]string, 0)
	err := c.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}
	w := &kafka.Writer{
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.LeastBytes{},
	}
	return &kafkaPublisher{writer: w, idg: idg}, nil
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
