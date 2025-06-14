package queue

import (
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

func NewKafkaProducer(config config.ApplicationConfig) *KafkaProducer {
	return &KafkaProducer{
		producer: *kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{config.KafkaWriterConfig.BrokerHost + ":" + config.KafkaWriterConfig.BrokerPort},
			Topic:   config.KafkaWriterConfig.Topic,
		}),
	}
}
