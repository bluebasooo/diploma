package queue

import (
	"dev/bluebasooo/video-recomendator/config"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	consumer kafka.Reader
}

func NewKafkaReader(config config.ApplicationConfig) *kafka.Reader {

}
