package kafka_sa

import (
	"context"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/google/wire"

	"purchase/infra/logx"
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
	writer    sarama.SyncProducer
	marshaler MessageMarshaler
}

func (s *kafkaPublisher) Publish(ctx context.Context, msg *mq.Message) error {
	msg.SetDeliveryTime(time.Now())
	kmsg, err := s.marshaler.Marshal(ctx, msg)
	if err != nil {
		return err
	}
	logx.Info(ctx, "publish message", slog.String("bizCode", msg.BizCode()), slog.String("topic", kmsg.Topic))
	_, _, err = s.writer.SendMessage(kmsg)
	return err
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
	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 3                    // 最大尝试发送次数
	conf.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	conf.Producer.Return.Successes = true

	// 连接 Kafka 服务器
	producer, err := sarama.NewSyncProducer(address, conf)
	if err != nil {
		return nil, err
	}
	return &kafkaPublisher{
		writer:    producer,
		marshaler: NewMessageMarshaler(idg),
	}, nil
}
