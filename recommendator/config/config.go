package config

import (
	"os"

	"dev/bluebasooo/video-common/config"
)

func GetApplicationConfig() *ApplicationConfig {
	return &ApplicationConfig{
		ClickhouseConfig: ClickhouseConfig{
			URI:      os.Getenv("CLICKHOUSE_URI"),
			User:     os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
			Database: os.Getenv("CLICKHOUSE_DATABASE"),
		},
		KafkaWriterConfig: KafkaWriterConfig{
			BrokerPort: os.Getenv("KAFKA_BROKER_URL"),
			BrokerHost: os.Getenv("KAFKA_BROKER_URL"),
			Topic:      os.Getenv("KAFKA_TOPIC"),
		},
		KafkaReaderConfig: KafkaReaderConfig{
			BrokerHost: os.Getenv("KAFKA_BROKER_URL"),
			BrokerPort: os.Getenv("KAFKA_BROKER_PORT"),
			Topic:      os.Getenv("KAFKA_TOPIC"),
			GroupID:    os.Getenv("KAFKA_GROUP_ID"),
		},
	}
}

type ApplicationConfig struct {
	ClickhouseConfig  ClickhouseConfig
	KafkaReaderConfig KafkaReaderConfig
	KafkaWriterConfig KafkaWriterConfig
	MongoConfig       config.MongoConfig
}

type KafkaReaderConfig struct {
	BrokerHost string
	BrokerPort string
	Topic      string
	GroupID    string
}

type KafkaWriterConfig struct {
	BrokerHost string
	BrokerPort string
	Topic      string
}

type ClickhouseConfig struct {
	URI      string
	User     string
	Password string
	Database string
}
