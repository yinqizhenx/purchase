package kafka

import (
	"context"

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
	// topic  string
	idg mq.IDGenFunc
}

func (s *kafkaPublisher) Publish(ctx context.Context, msg *mq.Message) error {
	kMsg := &kafka.Message{
		Key:   []byte(msg.HeaderGet(mq.EntityID)),
		Value: msg.Body,
	}
	m := NewMessage(kMsg)
	m.SetEventName(msg.HeaderGet(mq.EventName))
	m.SetMessageID(s.idg())

	kMsg, err := m.ToKafkaMessage()
	kMsg.Topic = msg.HeaderGet(mq.EventName)
	if err != nil {
		return err
	}
	err = s.writer.WriteMessages(ctx, *kMsg)
	if err != nil {
		return err
	}
	return nil
}

func (s *kafkaPublisher) Close() error {
	err := s.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewKafkaPublisher(c config.Config, idg mq.IDGenFunc) (mq.Publisher, error) {
	address := make([]string, 0)
	err := c.Value("kafka.address").Scan(&address)
	if err != nil {
		return nil, err
	}
	// topic, err := c.Value("kafka.topic").String()
	// if err != nil {
	// 	return nil, err
	// }
	w := &kafka.Writer{
		// Topic:    topic,
		Addr:     kafka.TCP(address...),
		Balancer: &kafka.LeastBytes{},
	}
	return &kafkaPublisher{writer: w, idg: idg}, nil
}
