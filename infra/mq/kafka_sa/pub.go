package kafka_sa

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
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
	writer sarama.SyncProducer
	idg    mq.IDGenFunc
}

func (s *kafkaPublisher) Publish(ctx context.Context, msg *mq.Message) error {
	msg.SetDeliveryTime(time.Now())
	kmsg, err := s.buildKafkaMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("发送消息bizCode:%s到topic：%s", msg.BizCode(), kmsg.Topic)
	_, _, err = s.writer.SendMessage(kmsg)
	return err
}

func (s *kafkaPublisher) buildKafkaMessage(m *mq.Message) (*sarama.ProducerMessage, error) {
	header := make([]sarama.RecordHeader, 0)
	header = append(header, sarama.RecordHeader{
		Key:   []byte(mq.MessageID),
		Value: []byte(s.idg()),
	})

	p, err := json.Marshal(m.Header())
	if err != nil {
		return nil, err
	}
	header = append(header, sarama.RecordHeader{
		Key:   []byte(HeaderPropertyKey),
		Value: p,
	})

	kmsg := &sarama.ProducerMessage{
		Key:     sarama.StringEncoder(m.BizCode()),
		Value:   sarama.ByteEncoder(m.Body),
		Topic:   m.EventName(),
		Headers: header,
	}

	// 重投的消息
	if m.RetryTopic() != "" {
		kmsg.Topic = m.RetryTopic()
	}

	// 死信消息
	if m.DeadTopic() != "" {
		kmsg.Topic = m.DeadTopic()
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
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 3                    // 最大尝试发送次数
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	config.Producer.Return.Successes = true

	// 连接 Kafka 服务器
	producer, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		return nil, err
	}
	return &kafkaPublisher{writer: producer, idg: idg}, nil
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
