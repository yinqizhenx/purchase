package kafka_sa

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct {
	ctx context.Context
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("topic=%s, partition=%d, offset=%d\n", msg.Topic, msg.Partition, msg.Offset)
		fmt.Printf("consume msg: %s\n", msg.Value)
		// 处理消息具体逻辑
	}
	return nil
}

type ConsumerGroup struct {
	cg sarama.ConsumerGroup
}

func NewConsumerGroup() (*ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true               // 成功发送的消息将写到 Successes 通道
	config.Consumer.Return.Errors = true                  // 消费时错误信息将写到 Errors 通道
	config.Consumer.Fetch.Default = 3 * 1024 * 1024       // 默认请求的字节数
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的 offset 读取，如果设置为 OffsetOldest 则从最旧的 offset 读取
	config.Consumer.Offsets.AutoCommit.Enable = true      // 将已消费的 offset 发送给 broker，默认为 true

	// 连接 Kafka 服务器，可以传入多个 broker，用逗号连接
	consumer, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "consumer-group", config)
	if err != nil {
		return nil, err
	}
	return &ConsumerGroup{cg: consumer}, nil
}

func (c ConsumerGroup) Consume(ctx context.Context) {

	// 消费消息
	for {
		err := c.cg.Consume(ctx, []string{"topic"}, &consumerGroupHandler{}) // 传入定义好的 ConsumerGroupHandler 结构体
		if err != nil {
			fmt.Printf("consume error: %#v\n", err)
		}
	}
}

func (c ConsumerGroup) Close() error {
	return c.cg.Close()
}
