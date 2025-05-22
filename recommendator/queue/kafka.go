package queue

import (
	"context"
	"dev/bluebasooo/video-recomendator/config"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	consumer kafka.Reader
}

type KafkaProducer struct {
	producer kafka.Writer
}

func NewKafkaConsumer(config config.ApplicationConfig) *KafkaConsumer {
	return &KafkaConsumer{
		consumer: *kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{config.KafkaReaderConfig.BrokerHost + ":" + config.KafkaReaderConfig.BrokerPort},
			Topic:          config.KafkaReaderConfig.Topic,
			GroupID:        config.KafkaReaderConfig.GroupID,
			CommitInterval: 0,
		}),
	}
}

func (k *KafkaConsumer) Consume(ctx context.Context, handler func(ctx context.Context, message *kafka.Message)) error {
	for {
		message, err := k.consumer.ReadMessage(ctx)
		if err != nil {
			return err
		}
		handler(ctx, &message)

		// TODO store task and commit after execution with wait groups
		err = k.consumer.CommitMessages(ctx, message)
		if err != nil {
			return err
		}
	}
}

func NewKafkaProducer(config config.ApplicationConfig) *KafkaProducer {
	return &KafkaProducer{
		producer: *kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{config.KafkaWriterConfig.BrokerHost + ":" + config.KafkaWriterConfig.BrokerPort},
			Topic:   config.KafkaWriterConfig.Topic,
		}),
	}
}

func (k *KafkaProducer) Produce(ctx context.Context, message *kafka.Message) error {
	return k.producer.WriteMessages(ctx, *message)
}
